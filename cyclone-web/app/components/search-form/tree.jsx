import React from 'react';
import { Tree } from 'antd';
const TreeNode = Tree.TreeNode;

export default class searchTree extends React.Component {

  onSelect = (selectedKeys, e) => {
    const data = e.selectedNodes[0].props.data;
    this.props.onSelect({ id: data.id, name: data.name });
  };

  render() {
    const treeData = this.props.treeData;
    const nodes = [treeData];
    if (nodes.length === 0) {
      return (
        <div style={{ color: '#999' }}>
          暂无数据
        </div>
      );
    }
    return (
      <Tree
        onSelect={this.onSelect}
        defaultExpandAll={true}
        defaultExpandParent={true}
      >
        {this.renderTreeNodes(nodes)}
      </Tree>
    );
  }

  renderTreeNodes = data => {
    return data.map(item => {
      if (item.children && item.children.length > 0) {
        return (
          <TreeNode data={item} key={item.id} title={item.name} onClick={(item) => this.onSelect(item)}>
            {this.renderTreeNodes(item.children)}
          </TreeNode>
        );
      }
      return <TreeNode data={item} key={item.id} title={item.name} onClick={(item) => this.onSelect(item)}/>;
    });
  };
}
