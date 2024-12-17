/**
 * 杭州云霁科技有限公司
 * http://www.idcos.com
 * Copyright (c) 2016 All Rights Reserved.
 */
package com.idcos.enterprise.portal.vo;

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
@SuppressWarnings("serial")
public class Item implements Serializable {

    private String id;

    private String title;

    private List<Item> children = new ArrayList<>();

    @JsonIgnore
    private Map<String, Item> childrenUnique = new HashMap<>();

    public Item() {

    }

    public Item(String id, String title) {
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

    public List<Item> getChildren() {
        return children;
    }

    public void setChildren(List<Item> children) {
        this.children = children;
    }

    /**
     * return child in this
     *
     * @param item
     * @return
     */
    public Item addChild(Item item) {
        if (childrenUnique.containsKey(item.getId())) {
            return childrenUnique.get(item.getId());
        }
        childrenUnique.put(item.getId(), item);
        this.children.add(item);
        return item;
    }

}
