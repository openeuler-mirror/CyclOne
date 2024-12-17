package com.idcos.enterprise.portal.biz.common.utils;

import com.idcos.cloud.biz.common.check.CommonParamtersChecker;
import com.idcos.cloud.core.common.biz.CommonResult;
import com.idcos.cloud.core.common.util.ListUtil;
import com.idcos.cloud.core.dal.common.page.PageForm;
import com.idcos.cloud.core.dal.common.page.PageUtils;
import com.idcos.cloud.core.dal.common.page.Pagination;
import com.idcos.cloud.core.dal.common.query.OperatorEnum;
import com.idcos.cloud.core.dal.common.query.SearchConditionBuilder;
import com.idcos.enterprise.portal.biz.common.CommonBizException;
import com.idcos.enterprise.portal.biz.common.ResultCode;
import com.idcos.enterprise.portal.biz.common.convert.BaseConvertFunction;
import com.idcos.enterprise.portal.biz.common.tempalte.*;
import com.idcos.enterprise.portal.dal.enums.IsActiveEnum;
import com.idcos.enterprise.portal.dal.repository.PortalGroupRoleRelRepository;
import com.idcos.enterprise.portal.dal.repository.PortalGroupUserRelRepository;
import com.idcos.enterprise.portal.dal.repository.PortalTenantRepository;
import org.apache.commons.beanutils.PropertyUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.support.ApplicationObjectSupport;
import org.springframework.data.domain.Page;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.stereotype.Service;

import java.lang.reflect.InvocationTargetException;
import java.util.List;

/**
 * 基本的查询删除功能实现
 *
 * @author jiaohuizhe
 * @version $Id: CrudUtilService.java, v 0.1 2015年5月27日 下午2:24:05 jiaohuizhe Exp $
 */
@Service
public class CrudUtilService extends ApplicationObjectSupport {
    protected final Logger logger = LoggerFactory.getLogger(getClass());
    @Autowired
    private BusinessProcessTemplate businessProcessTemplate;

    @Autowired
    private BusinessQueryTemplate businessQueryTemplate;

    @Autowired
    private PortalTenantRepository portalTenantRepository;

    @Autowired
    private PortalGroupUserRelRepository portalGroupUserRelRepository;

    @Autowired
    private PortalGroupRoleRelRepository portalGroupRoleRelRepository;

    /**
     * 根据数据ID和转换类删除某条数据
     *
     * @param id   数据ID
     * @param func PO2VO转换类 PO持久类类型 VO界面展示类类型
     * @return
     */
    public <PO, VO> CommonResult<VO> delete(final String id, final BaseConvertFunction<PO, VO> func) {
        return businessProcessTemplate.process(new BusinessProcessCallback<VO>() {
            @Override
            public void checkParam(BusinessProcessContext context) {
                CommonParamtersChecker.checkNotBlank(id, "删除【" + func.getTabEnum().getDescription()
                        + "】的功能接收参数ID不能为空！");
            }

            @Override
            public void checkBusinessInfo(BusinessProcessContext context) {
            }

            @Override
            @SuppressWarnings("unchecked")
            public VO doBusiness(BusinessProcessContext context) {
                JpaRepository<PO, String> repository = getApplicationContext().getBean(
                        func.getTabEnum().getBeanName(), JpaRepository.class);
                PO po = repository.findOne(id);
                if (po == null) {
                    throw new CommonBizException(ResultCode.QUERY_RESULT_IS_NULL, "系统未查询到要删除的数据！");
                } else {
                    try {
                        if (IsActiveEnum.HAS_ACTIVE.getCode().equals(
                                PropertyUtils.getSimpleProperty(po, "isActive"))) {
                            PropertyUtils.setSimpleProperty(po, "isActive",
                                    IsActiveEnum.NO_ACTIVE.getCode());
                            return func.apply(repository.save(po));
                        } else {
                            throw new CommonBizException(ResultCode.PARAM_ERROR,
                                    "当前数据状态为已删除，无法再次删除！");
                        }
                    } catch (ReflectiveOperationException e) {
                        throw new CommonBizException(ResultCode.UNKNOWN_EXCEPTION, "删除数据时出现未知错误！",
                                e);
                    }
                }
            }

            @Override
            public void exceptionProcess(CommonBizException exception,
                                         BusinessProcessContext context) {
            }
        });
    }


    /**
     * 根据form提交参数中的注解查询数据
     *
     * @param pclass   持久类类型
     * @param form     前端提交表单内容
     * @param pageForm 提交提交分页内容
     * @param func     PO2VO转换类
     * @return
     */
    public <PO, VO> CommonResult<Pagination<VO>> query(final Class<PO> pclass, final Object form,
                                                       final PageForm pageForm,
                                                       final BaseConvertFunction<PO, VO> func,
                                                       final boolean isActive) {
        CommonResult<Pagination<VO>> result = businessQueryTemplate
                .process(new BusinessQueryCallback<Pagination<VO>>() {
                    @Override
                    public void checkParam() {
                    }

                    @Override
                    @SuppressWarnings("unchecked")
                    public Pagination<VO> doQuery() {
                        JpaSpecificationExecutor<PO> repository = getApplicationContext().getBean(
                                func.getTabEnum().getBeanName(), JpaSpecificationExecutor.class);
                        SearchConditionBuilder<PO> scb = SearchConditionBuilder.newInstance(pclass);
                        if (isActive) {
                            scb.builder("isActive", OperatorEnum.EQ, IsActiveEnum.HAS_ACTIVE.getCode());
                        }
                        Page<PO> page = repository.findAll(scb.builder(form).createQuery(),
                                PageUtils.buildPageRequest(pageForm));
                        return PageUtils.toPagination(page, func);
                    }
                });
        return result;
    }

    /**
     * 根据form提交参数中的注解查询数据,此方法中不做分页处理,仅处理pageForm中的排序信息.
     *
     * @param pclass   持久类类型
     * @param form     前端提交表单内容
     * @param pageForm 提交提交排序内容
     * @param func     PO2VO转换类
     * @return
     */
    public <PO, VO> CommonResult<List<VO>> queryAll(final Class<PO> pclass, final Object form,
                                                    final PageForm pageForm,
                                                    final BaseConvertFunction<PO, VO> func,
                                                    final boolean isActive) {
        CommonResult<List<VO>> result = businessQueryTemplate
                .process(new BusinessQueryCallback<List<VO>>() {
                    @Override
                    public void checkParam() {
                    }

                    @Override
                    @SuppressWarnings("unchecked")
                    public List<VO> doQuery() {
                        JpaSpecificationExecutor<PO> repository = getApplicationContext().getBean(
                                func.getTabEnum().getBeanName(), JpaSpecificationExecutor.class);
                        SearchConditionBuilder<PO> scb = SearchConditionBuilder.newInstance(pclass);
                        if (isActive) {
                            scb.builder("isActive", OperatorEnum.EQ, IsActiveEnum.HAS_ACTIVE.getCode());
                        }
                        List<PO> list = repository.findAll(scb.builder(form).createQuery(), PageUtils
                                .buildPageRequest(pageForm).getSort());
                        return ListUtil.transform(list, func);
                    }
                });
        return result;
    }

    /**
     * 根据数据ID和转换类查询某条数据
     *
     * @param id   数据ID
     * @param func PO2VO转换类 PO持久类类型 VO界面展示类类型
     * @return
     */
    public <PO, VO> CommonResult<VO> findOne(final String id, final BaseConvertFunction<PO, VO> func) {
        CommonResult<VO> result = businessQueryTemplate.process(new BusinessQueryCallback<VO>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            @SuppressWarnings("unchecked")
            public VO doQuery() {
                JpaRepository<PO, String> repository = getApplicationContext().getBean(
                        func.getTabEnum().getBeanName(), JpaRepository.class);
                PO po = repository.findOne(id);
                if (po == null) {
                    throw new CommonBizException(ResultCode.QUERY_RESULT_IS_NULL, "查询结果为空！");
                } else {
                    return func.apply(po);
                }
            }
        });
        return result;
    }

    /**
     * 根据数据ID和转换类查询某条数据，并验证有效性
     *
     * @param id   数据ID
     * @param func PO2VO转换类 PO持久类类型 VO界面展示类类型
     * @return
     */
    public <PO, VO> CommonResult<VO> findOneAndIsActive(final String id,
                                                        final BaseConvertFunction<PO, VO> func) {
        CommonResult<VO> result = businessQueryTemplate.process(new BusinessQueryCallback<VO>() {
            @Override
            public void checkParam() {
                CommonParamtersChecker.checkNotBlank(id);
            }

            @Override
            @SuppressWarnings("unchecked")
            public VO doQuery() {
                JpaRepository<PO, String> repository = getApplicationContext().getBean(
                        func.getTabEnum().getBeanName(), JpaRepository.class);
                PO po = repository.findOne(id);
                if (po == null) {
                    throw new CommonBizException(ResultCode.QUERY_RESULT_IS_NULL, "查询结果为空！");
                } else {
                    try {
                        if (!IsActiveEnum.HAS_ACTIVE.getCode().equals(
                                PropertyUtils.getSimpleProperty(po, "isActive"))) {
                            throw new CommonBizException(ResultCode.STATUS_ERROR,
                                    "当前" + func.getTabEnum().getDescription() + "已被删除！");
                        }
                    } catch (IllegalAccessException | InvocationTargetException
                            | NoSuchMethodException e) {
                        throw new CommonBizException(ResultCode.STATUS_ERROR, "查询异常", e);
                    }
                    return func.apply(po);
                }
            }
        });
        return result;
    }

}
