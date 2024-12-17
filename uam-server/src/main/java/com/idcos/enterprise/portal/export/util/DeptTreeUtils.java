/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2016 All Rights Reserved.
 */
package com.idcos.enterprise.portal.export.util;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import com.idcos.enterprise.portal.dal.entity.PortalDept;
import com.idcos.enterprise.portal.dal.enums.SourceTypeEnum;
import org.apache.commons.collections.CollectionUtils;

import java.util.Collections;
import java.util.HashMap;
import java.util.List;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * 辅助前端生成树形结构
 *
 * @author Dana
 * @version DeptTreeUtils.java, v1 2017/11/03 下午1:34 Dana Exp $$
 */
public class DeptTreeUtils {

    public static Object getTree(List<PortalDept> list, TreeStyle treeStyle) {

        checkNotNull(treeStyle);

        if (CollectionUtils.isEmpty(list)) {
            return Collections.emptyList();
        }

        HashMap<String, TreeNode> treeNodeMap = Maps.newHashMap();
        List<TreeNode> treeNodes = Lists.newLinkedList();

        for (PortalDept portalDept : list) {

            TreeNode child = new TreeNode(treeStyle);
            child.setId(portalDept.getId());
            child.setTitle(portalDept.getDisplayName());
            child.setSourceType(portalDept.getSourceType());

            treeNodeMap.put(portalDept.getId(), child);
        }

        for (PortalDept portalDept : list) {

            TreeNode parent = treeNodeMap.get(portalDept.getParentId());

            TreeNode treeNode = treeNodeMap.get(portalDept.getId());

            if (parent == null) {
                treeNodes.add(treeNode);
                continue;
            }

            parent.getChildren().add(treeNode);
        }

        return treeNodes;
    }

    public static final class TreeNode extends HashMap<String, Object> {

        private static final long serialVersionUID = 1L;

        @JsonIgnore
        private TreeStyle treeStyle;

        public TreeNode(TreeStyle treeStyle) {
            this.treeStyle = checkNotNull(treeStyle);

            this.put(treeStyle.getIdName(), "");
            this.put(treeStyle.getTitleName(), "");
            this.put(treeStyle.getSourceType(), "");
            this.put(treeStyle.getChildrenName(), Lists.newLinkedList());
        }

        public void setId(String id) {
            this.put(treeStyle.getIdName(), id);
        }

        public void setTitle(String title) {
            this.put(treeStyle.getTitleName(), title);
        }

        public void setSourceType(String sourceType) {
            this.put(treeStyle.getSourceType(), sourceType);
        }

        public List<TreeNode> getChildren() {
            return (List<TreeNode>) this.get(treeStyle.getChildrenName());
        }

    }

    public enum TreeStyle {
        /**
         * io_tree
         */
        IO_TREE("id", "title", "children", "_node", "sourceType");

        private final String idName;

        private final String titleName;

        private final String childrenName;

        private final String nodeName;

        private final String sourceType;

        TreeStyle(String idName, String titleName, String childrenName, String nodeName, String sourceType) {
            this.idName = checkNotNull(idName);
            this.titleName = checkNotNull(titleName);
            this.childrenName = checkNotNull(childrenName);
            this.nodeName = checkNotNull(nodeName);
            this.sourceType = checkNotNull(sourceType);
        }

        public String getIdName() {
            return idName;
        }

        public String getTitleName() {
            return titleName;
        }

        public String getChildrenName() {
            return childrenName;
        }

        public String getNodeName() {
            return nodeName;
        }

        public String getSourceType() {
            return sourceType;
        }
    }

}
