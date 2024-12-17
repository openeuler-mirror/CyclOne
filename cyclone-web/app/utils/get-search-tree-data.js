import React from 'react';

/**
 * return treeData
 */
const isLeaf = (node) => {
  //return node.leaf || !node.children;
  return node.children.length === 0;
};

const searchAllNodes = (nodes, search) => {
  const expandedKeys = [];
  const walk = (nodes = []) => {
    let find = false;
    if (nodes.length === 0) {
      return find;
    }
    nodes.forEach(node => {
      // reset
      node.find = false;
      if (node.children && node.children.length > 0) {
        node.open = walk(node.children);
        if (node.open) {
          find = true;
          expandedKeys.push(node.id);
        }
      }
      if (node.name.indexOf(search) >= 0) {
        find = true;
        node.find = true;
      }
    });
    return find;
  };

  walk(nodes);
  return {
    expandedKeys,
    nodes
  };
};


const searchLeafOnly = (nodes, search) => {
  const ret = [];
  const walk = (nodes = []) => {
    if (nodes.length === 0) {
      return;
    }
    nodes.forEach(node => {
      if (node.children && node.children.length > 0) {
        walk(node.children);
      }
      if (isLeaf(node)) {
        if (node.name.indexOf(search) >= 0) {
          ret.push(node);
        }
      }
    });
  };

  walk(nodes);
  return ret;
};

const getSearchTreeData = (nodes = [], search = '', options = {}) => {
  const {
    searchDirectory = false
    // labelField = "name"
  } = options;

  if (!search || search.trim() === '') {
    if (!searchDirectory) {
      return nodes;
    } else {
      return {
        nodes,
        expandedKeys: []
      };
    }
  }

  search = search.trim();

  if (!searchDirectory) {
    return searchLeafOnly(nodes, search);
  }

  return searchAllNodes(nodes, search);
};

export const getSearchLabel = (text, searchValue) => {
  searchValue = searchValue.trim();
  const index = text.indexOf(searchValue);
  if (searchValue === '' || index < 0) {
    return text;
  }
  const pre = text.substring(0, index);
  const after = text.substring(index + searchValue.length);
  return (
    <span>
      {pre}
      <b style={{
        color: '#000000',
        background: '#fbdd17',
        display: 'inline-block',
        padding: '0 5px'
      }}
      > {searchValue} </b>
      {after}
    </span>
  );
};

export default getSearchTreeData;
