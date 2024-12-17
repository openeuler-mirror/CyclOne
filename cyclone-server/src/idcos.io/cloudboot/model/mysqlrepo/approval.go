package mysqlrepo

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/voidint/page"

	"idcos.io/cloudboot/model"
)

// SubmitApproval 提交审批单及其审批步骤
func (repo *MySQLRepo) SubmitApproval(approval *model.Approval, steps ...*model.ApprovalStep) (err error) {
	tx := repo.db.Begin()

	if err = tx.Create(approval).Error; err != nil {
		tx.Rollback()
		repo.log.Error(err)
		return err
	}
	for i := range steps {
		steps[i].ApprovalID = approval.ID
		if err = tx.Create(steps[i]).Error; err != nil {
			tx.Rollback()
			repo.log.Error(err)
			return err
		}
	}
	return tx.Commit().Error
}

// Approve 审批
func (repo *MySQLRepo) Approve(approvalID, stepID uint, action, remark string) error {
	// TODO 待实现
	return nil
}

// UpdateApproval 修改审批单
func (repo *MySQLRepo) UpdateApproval(mod *model.Approval) (affected int64, err error) {
	err = repo.db.Model(&model.Approval{}).Where("id = ?", mod.ID).Update(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// RevokeApproval 撤销目标审批单
func (repo *MySQLRepo) RevokeApproval(approvalID uint) (err error) {
	now := time.Now()
	if err = repo.db.Model(&model.Approval{}).Where("id = ? and status = ?", approvalID, model.ApprovalStatusApproval).Updates(&model.Approval{
		Status:  model.ApprovalStatusRevoked,
		EndTime: &now,
	}).Error; err != nil {
		repo.log.Error(err)
		return err
	}
	return nil
}

// GetApprovalByID 查询指定审批单
func (repo *MySQLRepo) GetApprovalByID(approvalID uint) (approval *model.Approval, err error) {
	var ap model.Approval
	if err := repo.db.Where("id = ?", approvalID).Find(&ap).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return &ap, nil
}

// GetApprovalStepByApprovalID 查询指定审批单的审批步骤明细
func (repo *MySQLRepo) GetApprovalStepByApprovalID(approvalID uint) (steps []*model.ApprovalStep, err error) {
	steps = make([]*model.ApprovalStep, 0)
	db := repo.db.Model(&model.ApprovalStep{})
	err = db.Where("approval_id = ?", approvalID).Find(&steps).Error
	return
}

// CountPendingApprovals 统计审批单个数
func (repo *MySQLRepo) CountPendingApprovals(currentUserID string, cond *model.Approval) (count int64, err error) {
	steps, err := repo.GetApprovalStepByCond(&model.ApprovalStep{
		Approver: currentUserID,
	}, false)
	if err != nil {
		return 0, err
	}
	var ids []uint
	for k := range steps {
		ids = append(ids, steps[k].ApprovalID)
	}

	// 待我审批肯定在step表中，如果没有则不存在待我审批的审批单
	if len(ids) == 0 {
		return 0, nil
	}

	// 根据审批流程里获取的审批单ID 获取审批单
	db := repo.db.Model(&model.Approval{})
	if cond != nil {
		if cond.Type != "" {
			db = db.Where("type = ?", cond.Type)
		}
	}

	db = db.Where("status = ?", model.ApprovalStatusApproval)

	if len(ids) > 0 {
		db = db.Where("id in (?)", ids)
	}

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// CountInitiatedApprovals 统计审批单个数
func (repo *MySQLRepo) CountInitiatedApprovals(currentUserID string, cond *model.Approval) (count int64, err error) {
	db := repo.db.Model(&model.Approval{})
	if cond != nil {
		if cond.Type != "" {
			db = db.Where("type = ?", cond.Type)
		}
		if cond.Status != "" {
			db = db.Where("status = ?", cond.Status)
		}
		if cond.ID > 0 {
			db = db.Where("id = ?", cond.ID)
		}
	}
	db = db.Where("initiator = ?", currentUserID)

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetInitiatedApprovals 查询'我发起的'审批分页列表
func (repo *MySQLRepo) GetInitiatedApprovals(currentUserID string, cond *model.Approval, orderby model.OrderBy, limiter *page.Limiter) (items []*model.Approval, err error) {
	db := repo.db.Model(&model.Approval{})
	if cond != nil {
		if cond.Type != "" {
			db = db.Where("type = ?", cond.Type)
		}
		if cond.Status != "" {
			db = db.Where("status = ?", cond.Status)
		}
		if cond.ID > 0 {
			db = db.Where("id = ?", cond.ID)
		}
	}
	db = db.Where("initiator = ?", currentUserID)

	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// GetPendingApprovals 查询'待我审批'的审批分页列表
func (repo *MySQLRepo) GetPendingApprovals(currentUserID string, cond *model.Approval, orderby model.OrderBy, limiter *page.Limiter) (items []*model.Approval, err error) {
	// 在审批流程中查询被currentUserID审批的审批单
	steps, err := repo.GetApprovalStepByCond(&model.ApprovalStep{
		Approver: currentUserID,
	}, false)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for k := range steps {
		ids = append(ids, steps[k].ApprovalID)
	}

	// 待我审批肯定在step表中，如果没有则不存在待我审批的审批单
	if len(ids) == 0 {
		return nil, nil
	}

	// 根据审批流程里获取的审批单ID 获取审批单
	db := repo.db.Model(&model.Approval{})
	if cond != nil {
		if cond.Type != "" {
			db = db.Where("type = ?", cond.Type)
		}
	}

	db = db.Where("status = ?", model.ApprovalStatusApproval)

	if len(ids) > 0 {
		db = db.Where("id in (?)", ids)
	}

	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// CountApprovedApprovals 统计'我已审批的'的审批表个数
func (repo *MySQLRepo) CountApprovedApprovals(currentUserID string, cond *model.Approval) (count int64, err error) {
	// 在审批流程中查询被currentUserID审批的审批单
	steps, err := repo.GetApprovalStepByCond(&model.ApprovalStep{
		Approver: currentUserID,
	}, true)
	if err != nil {
		return count, err
	}
	var ids []uint
	for k := range steps {
		ids = append(ids, steps[k].ApprovalID)
	}

	// 审批步骤中没有，则不存在被我审批过的审批单
	if len(ids) == 0 {
		return count, nil
	}

	db := repo.db.Model(&model.Approval{})
	if cond != nil {
		if cond.Type != "" {
			db = db.Where("type = ?", cond.Type)
		}
		if cond.Status != "" {
			db = db.Where("status = ?", cond.Status)
		}
	}

	if len(ids) > 0 {
		db = db.Where("id in (?)", ids)
	}

	if err = db.Select("id").Count(&count).Error; err != nil {
		repo.log.Error(err)
		return count, err
	}
	return count, nil
}

// GetApprovedApprovals 查询'我已审批的'的审批分页列表
func (repo *MySQLRepo) GetApprovedApprovals(currentUserID string, cond *model.Approval, orderby model.OrderBy, limiter *page.Limiter) (items []*model.Approval, err error) {
	// 在审批流程中查询被currentUserID审批的审批单
	steps, err := repo.GetApprovalStepByCond(&model.ApprovalStep{
		Approver: currentUserID,
	}, true)
	if err != nil {
		return nil, err
	}
	var ids []uint
	for k := range steps {
		ids = append(ids, steps[k].ApprovalID)
	}

	// 审批步骤中没有，则不存在被我审批过的审批单
	if len(ids) == 0 {
		return nil, nil
	}

	db := repo.db.Model(&model.Approval{})
	if cond != nil {
		if cond.Type != "" {
			db = db.Where("type = ?", cond.Type)
		}
		if cond.Status != "" {
			db = db.Where("status = ?", cond.Status)
		}
	}
	if len(ids) > 0 {
		db = db.Where("id in (?)", ids)
	}

	for i := range orderby {
		db = db.Order(orderby[i].String())
	}
	if limiter != nil {
		db = db.Limit(limiter.Limit).Offset(limiter.Offset)
	}

	if err = db.Find(&items).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return items, nil
}

// GetApprovalStepByCond 根据条件查询审批单的审批步骤明细
func (repo *MySQLRepo) GetApprovalStepByCond(cond *model.ApprovalStep, isApproved bool) (steps []*model.ApprovalStep, err error) {
	db := repo.db.Model(&model.ApprovalStep{})
	if cond != nil {
		if cond.ApprovalID > 0 {
			db = db.Where("approval_id = ?", cond.ApprovalID)
		}
		// 一个步骤经过多个审批人同时审批,这里需要修改
		if cond.Approver != "" {
			db = db.Where("approver = ?", cond.Approver)
		}
	}

	if !isApproved {
		// 如果是待审批，end_time字段应该为空
		db = db.Where("end_time IS NULL")
		db = db.Where("start_time IS NOT NULL")
		db = db.Where("start_time < ?", time.Now())
	} else {
		// 如果是已经审批，end_time字段应该是现在以前的时间
		db = db.Where("end_time < ?", time.Now())
		db = db.Where("start_time < ?", time.Now())
		db = db.Where("start_time IS NOT NULL")
		db = db.Where("end_time IS NOT NULL")
	}

	if err = db.Find(&steps).Error; err != nil {
		repo.log.Error(err)
		return nil, err
	}
	return steps, nil
}

//GetApprovalStepByID 根据审批步骤ID查询审批步骤
func (repo *MySQLRepo) GetApprovalStepByID(stepID uint) (step *model.ApprovalStep, err error) {
	step = &model.ApprovalStep{}
	if err := repo.db.Where("id = ?", stepID).Find(step).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			repo.log.Error(err)
		}
		return nil, err
	}
	return step, nil
}

// UpdateApprovalStep 修改审批步骤
func (repo *MySQLRepo) UpdateApprovalStep(mod *model.ApprovalStep) (affected int64, err error) {
	err = repo.db.Model(&model.ApprovalStep{}).Where("id = ?", mod.ID).Updates(mod).Error
	if err != nil {
		return 0, err
	}
	return 1, nil
}
