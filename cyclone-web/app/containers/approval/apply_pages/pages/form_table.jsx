import React from 'react';
import { Table } from 'antd';

export default class DeviceTable extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      dataSource: props.dataSource || []
    };
  }

  componentWillReceiveProps(props) {
    if (props.dataSource.length > 0) {
      this.setState({
        dataSource: props.dataSource
      });
    }
  }

  handleDelete = (id) => {
    const dataSource = [...this.state.dataSource];
    const newDataSource = dataSource.filter(item => item.id !== id);
    this.setState({ dataSource: newDataSource });
    this.props.setFormValue(newDataSource);
  };

  getLocalColumns = () => {
    let columns = [];
    const baseColumns = this.props.tableColumn();
    const operateColumns = [
      {
        title: '操作',
        dataIndex: 'operation',
        render: (text, record) => {
          return (
            this.state.dataSource.length >= 1
              ?
                <a href='javascript:;' onClick={() => this.handleDelete(record.id)} style={{ color: 'rgb(255, 55, 0)' }}>删除</a>
              : null
          );
        }
      }
    ];
    columns = [ ...baseColumns, ...operateColumns ];
    return columns;
  };

  render() {
    const { dataSource } = this.state;
    return (
      <Table
        bordered={true}
        dataSource={dataSource}
        columns={this.getLocalColumns()}
        pagination={false}
      />
    );
  }
}
