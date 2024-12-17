import React from 'react';
import { Table, Tooltip, notification, Badge, Alert, Select, Popconfirm } from 'antd';
import { OPERATION_STATUS_COLOR } from "common/enums";
import { getWithArgs } from 'common/xFetch2';

export default class DeviceTable extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      dataSource: props.dataSource || [],
      cabinet: [],
      usite: [],
      idcs: [],
      rooms: [],
      virtualCabinets: [],
      dst_server_cabinet: {},
      dst_server_room: {},
      dst_server_usite: {},
      dst_virtual_cabinets: {},
      store_rooms: []

    };
  }

  componentWillReceiveProps(props) {
    if (props.dataSource.length > 0) {
      this.setState({
        dataSource: props.dataSource
      });
    }
  }

  componentDidMount() {
    getWithArgs('/api/cloudboot/v1/idcs', { page: 1, page_size: 100 }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState({
        idcs: res.content.records
      });
    });


  }

  //清空机房信息
  clearRoom = (record) => {
    delete record.dst_server_room_id;
    delete record.dst_server_room_name;
    this.setState({
      [`${record.sn}-rooms`]: [],
      [`${record.sn}-dst_server_room`]: {}
    });
  };
  //清空机架信息
  clearCabinet = (record) => {
    delete record.dst_cabinet_id;
    delete record.dst_cabinet_number;
    this.setState({
      [`${record.sn}-cabinet`]: [],
      [`${record.sn}-dst_server_cabinet`]: {}
    });
  };
  //清空机位信息
  clearUsite = (record) => {
    delete record.dst_usite_id;
    delete record.dst_usite_number;
    this.setState({
      [`${record.sn}-usite`]: [],
      [`${record.sn}-dst_server_usite`]: {}
    });
  };

  //选择库房获取虚机货架
  getVirtualCabinets = (v, record, page) => {
    record.dst_store_room_id = v.key;
    record.dst_store_room_name = v.label;
    this.setState({
      [`${record.sn}-dst_store_room`]: {
        key: v.key,
        label: v.label
      }
    });
    this.getVirtualCabinetsPage(v, record, page);
  };
  getVirtualCabinetsPage = (v, record, page) => {
    getWithArgs('/api/cloudboot/v1/virtual-cabinets', { page: page, page_size: 100, store_room_id: v.key }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          [`${record.sn}-virtualCabinets`]: [ ...preState.virtualCabinets, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getVirtualCabinetsPage(v, record, page + 1);
        }
      });
    });
  };

  //选择数据中心获取机房列表
  getRooms = (v, record, page) => {
    record.dst_idc_id = v.key;
    record.dst_idc_name = v.label;
    this.setState({
      [`${record.sn}-dst_idc`]: {
        key: v.key,
        label: v.label
      }
    });
    this.clearRoom(record);
    this.clearCabinet(record);
    this.clearUsite(record);
    this.props.setFormValue(this.state.dataSource);
    this.getRoomPage(v, record, page);
    this.getStoreRoomPage(v, record, page);


  };

  getRoomPage = (v, record, page) => {
    getWithArgs('/api/cloudboot/v1/server-rooms', { page: page, page_size: 100, idc_id: v.key }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          [`${record.sn}-rooms`]: [ ...preState.rooms, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getRoomPage(v, record, page + 1);
        }
      });
    });
  };


  getStoreRoomPage = (v, record, page) => {
    getWithArgs('/api/cloudboot/v1/store-rooms', { page: page, page_size: 100, idc_id: v.key }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          [`${record.sn}-store_rooms`]: [ ...preState.store_rooms, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getStoreRoomPage(v, record, page + 1);
        }
      });
    });
  };


  //选择机房获取机架列表
  getCabinet = (v, record, page) => {
    record.dst_server_room_id = v.key;
    record.dst_server_room_name = v.label;
    this.setState({
      [`${record.sn}-dst_server_room`]: {
        key: v.key,
        label: v.label
      }
    });
    this.clearCabinet(record);
    this.clearUsite(record);
    this.props.setFormValue(this.state.dataSource);
    this.getCabinetPage(v, record, page);
  };
  getCabinetPage = (v, record, page) => {
    getWithArgs('/api/cloudboot/v1/server-cabinets', { page: page, page_size: 1000, server_room_id: v.key }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          [`${record.sn}-cabinet`]: [ ...preState.cabinet, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getCabinetPage(v, record, page + 1);
        }
      });
    });
  };

  //选择机架获取机位列表
  getUsite = (v, record, page) => {
    record.dst_cabinet_id = v.key;
    record.dst_cabinet_number = v.label;
    this.setState({
      [`${record.sn}-dst_server_cabinet`]: {
        key: v.key,
        label: v.label
      }
    });
    this.clearUsite(record);
    this.props.setFormValue(this.state.dataSource);
    this.getUsitePage(v, record, page);
  };
  getUsitePage = (v, record, page) => {
    getWithArgs('/api/cloudboot/v1/server-usites', { page: page, page_size: 100, server_cabinet_id: v.key, status: 'free,pre_occupied' }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          [`${record.sn}-usite`]: [ ...preState.usite, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getUsitePage(v, record, page + 1);
        }
      });
    });
  };

  setUsite = (v, record) => {
    record.dst_usite_id = v.key;
    record.dst_usite_number = v.label;
    this.setState({
      [`${record.sn}-dst_server_usite`]: {
        key: v.key,
        label: v.label
      }
    });
    this.props.setFormValue(this.state.dataSource);
  };

  setVirtualCabinets = (v, record) => {
    record.dst_virtual_cabinet_id = v.key;
    record.dst_virtual_cabinet_number = v.label;
    this.setState({
      [`${record.sn}-dst_virtual_cabinets`]: {
        key: v.key,
        label: v.label
      }
    });
    this.props.setFormValue(this.state.dataSource);
  };


  getColums = () => {
    return [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      }, {
        title: '序列号',
        dataIndex: 'sn',
        render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
      },
      {
        title: '设备类型',
        dataIndex: 'category',
        //render: (t) => {
        //  if (t === 'SpecialDev') {
        //    return '特殊设备';
        //  }
        //  return t;
        //}
      },
      {
        title: '运营状态',
        dataIndex: 'operation_status',
        render: type => {
          const color = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][0] : 'transparent';
          const word = OPERATION_STATUS_COLOR[type] ? OPERATION_STATUS_COLOR[type][1] : '';
          return (
            <div>
              <Badge
                dot={true}
                style={{
                  background: color
                }}
              />{' '}
              &nbsp;&nbsp; {word}
            </div>
          );
        }
      },
      {
        title: '原数据中心',
        dataIndex: 'idc',
        render: (t) => t ? t.name : ''
      },
      {
        title: '原机房',
        dataIndex: 'server_room',
        render: (t) => t ? t.name : ''
      },
      {
        title: '原机架/虚拟货架',
        dataIndex: 'server_cabinet',
        render: (t) => t ? t.number : ''
      },
      {
        title: '原机位',
        dataIndex: 'server_usite',
        render: (t) => t ? t.number : ''
      },
      {
        title: '目标数据中心',
        dataIndex: 'idc',
        width: 150,
        render: (text, record) => {
          return (
            <div>
              {/* <div className='move-table-seleted'> {record.dst_idc_name}</div>            */}
              <Select
                style={{ minWidth: 100 }}
                size='small'
                value={this.state[`${record.sn}-dst_idc`] || { key: record.dst_idc_id, label: record.dst_idc_name }}
                showSearch={true}
                filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                labelInValue={true} 
                onChange={(value) => this.getRooms(value, record, 1)}
              >
                {
                  (this.state.idcs || []).map((os, index) => <Option title={os.name} value={os.id} key={os.id}>{os.name}</Option>)
                }
              </Select>
            </div>

          );
        }
      }, {
        title: '目标机房',
        dataIndex: 'server_room',
        width: 150,
        required: true,
        render: (text, record) => {
          return (
            <div>
              {/* <div className='move-table-seleted'>{record.dst_server_room_name}</div>            */}
              <Select
                style={{ minWidth: 100 }}
                size='small'
                value={this.state[`${record.sn}-dst_server_room`] || { key: record.dst_server_room_id, label: record.dst_server_room_name }}
                showSearch={true}
                filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                labelInValue={true} onChange={(value) => this.getCabinet(value, record, 1)}
              >
                {
                  (this.state[`${record.sn}-rooms`] || []).map((os, index) => <Option title={os.name} value={os.id} key={os.id}>{os.name}</Option>)
                }
              </Select>
            </div>
          );
        }
      }, {
        title: '目标机架',
        dataIndex: 'server_cabinet',
        width: 150,
        required: true,
        render: (text, record) => {
          return (
            <div>
              {/* <div className='move-table-seleted'>{record.dst_cabinet_number}</div>             */}
              <Select
                style={{ minWidth: 100 }}
                size='small'
                value={this.state[`${record.sn}-dst_server_cabinet`] || { key: record.dst_cabinet_id, label: record.dst_cabinet_number }}
                showSearch={true}
                filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                labelInValue={true} onChange={(value) => this.getUsite(value, record, 1)}
              >
                {
                  (this.state[`${record.sn}-cabinet`] || []).map((os, index) => <Option title={os.number} value={os.id} key={os.id}>{os.number}</Option>)
                }
              </Select>
            </div>

          );
        }
      }, {
        title: '目标机位',
        dataIndex: 'server_usite',
        width: 150,
        render: (text, record) => {
          return (
            <div>
              {/* <div className='move-table-seleted'>{record.dst_usite_number}</div> */}
              <Select
                style={{ minWidth: 100 }}
                size='small'
                value={this.state[`${record.sn}-dst_server_usite`] || { key: record.dst_usite_id, label: record.dst_usite_number }}
                showSearch={true}
                filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                labelInValue={true} onChange={(value) => this.setUsite(value, record)}
              >
                {
                  (this.state[`${record.sn}-usite`] || []).map((os, index) => <Option title={os.number} value={os.id} key={os.id}>{os.number}</Option>)
                }
              </Select>
            </div>

          );
        }
      }, {
        title: '目标库房',
        dataIndex: 'dst_store_room',
        width: 150,
        render: (text, record) => {
          return (
            <div>
              {/* <div className='move-table-seleted'>{record.dst_store_room_name}</div> */}
              <Select
                style={{ minWidth: 100 }}
                value={this.state[`${record.sn}-dst_store_room`] || { key: record.dst_store_room_id, label: record.dst_store_room_name }}
                size='small'
                showSearch={true}
                filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                labelInValue={true} onChange={(value) => this.getVirtualCabinets(value, record, 1)}
              >
                {
                  (this.state[`${record.sn}-store_rooms`] || []).map((os, index) => <Option title={os.name} value={os.id} key={os.id}>{os.name}</Option>)
                }
              </Select>
            </div>
          );
        }
      }, {
        title: '目标虚拟货架',
        dataIndex: 'dst_virtual_cabinet',
        width: 150,
        render: (text, record) => {
          return (
            <div>
              {/* <div className='move-table-seleted'>{record.dst_virtual_cabinet_number}</div> */}
              <Select
                style={{ minWidth: 100 }}
                size='small'
                value={this.state[`${record.sn}-dst_virtual_cabinets`] || { key: record.dst_virtual_cabinet_id, label: record.dst_virtual_cabinet_number }}
                showSearch={true}
                filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                labelInValue={true} onChange={(value) => this.setVirtualCabinets(value, record)}
              >
                {
                  (this.state[`${record.sn}-virtualCabinets`] || []).map((os, index) => <Option value={os.id} title={os.number} key={os.id}>{os.number}</Option>)
                }
              </Select>
            </div>
          );
        }
      }, {
        title: '操作',
        dataIndex: 'operation',
        render: (text, record) => {
          return (
            this.state.dataSource.length >= 1
              ? (
                <Popconfirm title='确定删除吗？' onConfirm={() => this.handleDelete(record.id)}>
                  <a href='javascript:;' style={{ color: 'rgb(255, 55, 0)' }}>删除</a>
                </Popconfirm>
              ) : null
          );
        }
      }];
  };

  handleDelete = (id) => {
    const dataSource = [...this.state.dataSource];
    const newDataSource = dataSource.filter(item => item.id !== id);
    this.setState({ dataSource: newDataSource });
    this.props.setFormValue(newDataSource);
  };

  render() {
    const dataSource = this.state.dataSource;
    return (
      <div>
        <Alert closable={true} style={{ marginBottom: 8 }} message='机房和库房二选一，搬迁到机房或库房' type='info' showIcon={true} />
        <Table
          bordered={true}
          dataSource={dataSource}
          columns={this.getColums()}
          pagination={false}
        />

      </div>
     
    );
  }
}
