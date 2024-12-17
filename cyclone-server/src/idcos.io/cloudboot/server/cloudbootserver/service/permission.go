package service

import (
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// PermissionCodeNode 权限码节点
type PermissionCodeNode struct {
	ID       uint                  `json:"-"`
	PID      uint                  `json:"-"`
	Code     string                `json:"id"`
	Title    string                `json:"title"`
	Note     string                `json:"note"`
	Children []*PermissionCodeNode `json:"children"`
}

// PermissionCodeTree 查询权限码节点树
// 未检索到调用处，已废弃 - 2021-09-03
func PermissionCodeTree(log logger.Logger, repo model.Repo) (*PermissionCodeNode, error) {
	all, err := repo.GetPermissionCodes(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	nodes := make([]*PermissionCodeNode, 0, len(all))
	for _, node := range all {
		nodes = append(nodes, &PermissionCodeNode{
			ID:       node.ID,
			PID:      node.PID,
			Code:     node.Code,
			Title:    node.Title,
			Note:     node.Note,
			Children: []*PermissionCodeNode{},
		})
	}

	root := PermissionCodeNode{
		Title:    "所有权限",
		Children: []*PermissionCodeNode{},
	}

	for i := range nodes {
		var matched bool // 是否已经找到到父节点
		for j := range nodes {
			if nodes[i] == nil || nodes[j] == nil || nodes[i].ID == nodes[j].ID {
				continue
			}
			if nodes[i].PID == nodes[j].ID {
				nodes[j].Children = append(nodes[j].Children, nodes[i])
				matched = true
			}
		}
		if !matched {
			root.Children = append(root.Children, nodes[i])
		}
	}

	return &root, nil
}
