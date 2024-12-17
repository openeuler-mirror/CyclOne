/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2016 All Rights Reserved.
 */
package com.idcos.enterprise.portal.web.permission;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import com.fasterxml.jackson.annotation.JsonIgnore;

/**
 * @author zhouqin
 * @version com.idcos.automate.export.api.rbac.Item.java, v 1.1 5/12/16 zhouqin Exp $
 */
public class PermissionNode implements Serializable {

    /**  */
    private static final long serialVersionUID = 1L;

    private String id;

    private String title;

    private String note;

    private List<PermissionNode> children = new ArrayList<PermissionNode>();

    @JsonIgnore
    private Map<String, PermissionNode> childrenUnique = new HashMap<>();

    public PermissionNode() {

    }

    public PermissionNode(String id, String title) {
        this.id = id;
        this.title = title;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public List<PermissionNode> getChildren() {
        return children;
    }

    public void setChildren(List<PermissionNode> children) {
        this.children = children;
    }

    public String getNote() {
        return note;
    }

    public void setNote(String note) {
        this.note = note;
    }

    /**
     * return child in this
     *
     * @param permissionNode
     * @return
     */
    public PermissionNode addChild(PermissionNode permissionNode) {
        if (childrenUnique.containsKey(permissionNode.getId())) {
            return childrenUnique.get(permissionNode.getId());
        }
        childrenUnique.put(permissionNode.getId(), permissionNode);
        this.children.add(permissionNode);
        return permissionNode;
    }

}
