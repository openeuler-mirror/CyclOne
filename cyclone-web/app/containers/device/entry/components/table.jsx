import React from 'react';
import { Table, Input, Button, Popconfirm, Form, Modal, Badge } from 'antd';
const FormItem = Form.Item;
import actions from '../actions';
const EditableContext = React.createContext();
import { notification, Switch } from 'antd';
import { post } from 'common/xFetch2';
const confirm = Modal.confirm;
import { hashHistory } from 'react-router';
import { YES_NO, OPERATION_STATUS_COLOR } from "common/enums";
import ApprovalForm from 'containers/approval/apply_pages/reInstall/form';

const EditableRow = ({ form, index, ...props }) => (
  <EditableContext.Provider value={form}>
    <tr {...props} />
  </EditableContext.Provider>
);

const EditableFormRow = Form.create()(EditableRow);

class EditableCell extends React.Component {
  state = {
    editing: false
  };

  componentDidMount() {
    if (this.props.editable) {
      document.addEventListener('click', this.handleClickOutside, true);
    }
  }

  componentWillUnmount() {
    if (this.props.editable) {
      document.removeEventListener('click', this.handleClickOutside, true);
    }
  }

  toggleEdit = () => {
    const editing = !this.state.editing;
    this.setState({ editing: true }, () => {
      if (editing && this.input) {
        this.input.focus();
      }
    });
  };

  handleClickOutside = (e) => {
    const { editing } = this.state;
    if (editing && this.cell !== e.target && !this.cell.contains(e.target)) {
      this.save();
    }
  };

  save = () => {
    const { record, handleSave } = this.props;
    this.form.validateFields((error, values) => {
      if (error) {
        return;
      }
      this.toggleEdit();
      handleSave({ ...record, ...values });
    });
  };


  render() {
    const { editing } = this.state;
    const {
      editable,
      dataIndex,
      title,
      record,
      required,
      placeholder,
      ...restProps
    } = this.props;
    return (
      <td ref={node => (this.cell = node)} {...restProps}>
        {editable ? (
          <EditableContext.Consumer>
            {(form) => {
              this.form = form;
              return (
                editing ? (
                  <FormItem style={{ margin: 0 }}>
                    {form.getFieldDecorator(dataIndex, {
                      rules: [{
                        required: required,
                        message: `${title} 不能为空`
                      }],
                      initialValue: record[dataIndex]
                    })(
                      <Input
                        placeholder={placeholder}
                        ref={node => (this.input = node)}
                        onPressEnter={this.save}
                      />
                    )}
                  </FormItem>
                ) : (
                  <div
                    className='editable-cell-value-wrap'
                    style={{ paddingRight: 24 }}
                    onClick={this.toggleEdit}
                  >
                    {restProps.children}
                  </div>
                  )
              );
            }}
          </EditableContext.Consumer>
        ) : restProps.children}
      </td>
    );
  }
}

export default class EditableTable extends React.Component {
  constructor(props) {
    super(props);
    this.columns = [{
      title: '固资编号',
      dataIndex: 'fixed_asset_number',
      width: 150,
    }, {
      title: '序列号',
      dataIndex: 'sn',
      width: 150,
    }, {
      title: '操作系统',
      dataIndex: 'os_template_name',
      width: 200,
      required: true,
      render: (text, record) => {
        if (text) {
          return (
            <div style={{ position: 'relative' }}>
              <span className='hover-edit' onClick={() => this.execAction('addSystem', record)}>{text}</span>
            </div>
          );
        }
        return (
          <Button icon='plus' onClick={() => this.execAction('addSystem', record)}>操作系统</Button>
        );
      }
    }, {
      title: 'RAID类型',
      dataIndex: 'hardware_template_name',
      required: true,
      width: 200,
      render: (text, record) => {
        if (text) {
          return <span className='hover-edit' onClick={() => this.execAction('addHardware', record)}>{text}</span>;
        }
        return (
          <Button icon='plus' onClick={() => this.execAction('addHardware', record)}>RAID类型</Button>
        );
      }
    }, {
      title: '分配外网IPv4',
      dataIndex: 'need_extranet_ip',
      editable: false,
      required: true,
      width: 10,
      type: 'select',
      render: (text, record) => {
        return <Switch checkedChildren='是' unCheckedChildren='否' checked={text === 'yes'} onChange={(value) => this.onCheckNeedExIPChange(value, record)} />;
      }
    },
    {
      title: '分配内网IPv6',
      dataIndex: 'need_intranet_ipv6',
      width: 10,
      editable: false,
      required: true,
      type: 'select',
      render: (text, record) => {
        return <Switch checkedChildren='是' unCheckedChildren='否' checked={text === 'yes'} onChange={(value) => this.onCheckNeedInIPv6Change(value, record)} />;
      }
    }, 
    {
      title: '分配外网IPv6',
      dataIndex: 'need_extranet_ipv6',
      editable: false,
      required: true,
      width: 10,
      type: 'select',
      render: (text, record) => {
        return <Switch checkedChildren='是' unCheckedChildren='否' checked={text === 'yes'} onChange={(value) => this.onCheckNeedExIPv6Change(value, record)} />;
      }
    },        
    {
      title: '设备类型',
      dataIndex: 'category',
      width: 10,
    },
    {
      title: '用途',
      dataIndex: 'usage',
      width: 10,
    }, {
      title: '机架编号',
      dataIndex: 'server_cabinet',
      width: 10,
      render: (text) => <span>{text ? text.number : ''}</span>
    }, {
      title: '机位编号',
      dataIndex: 'server_usite',
      width: 10,
      render: (text) => <span>{text ? text.number : ''}</span>
    }, 
    {
      title: '物理区域',
      dataIndex: 'physical_area',
      width: 180,
      render: (t, record) => record.server_usite ? record.server_usite.physical_area : ''
    },
    {
      title: '运营状态',
      dataIndex: 'operation_status',
      width: 10,
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
            {word}
          </div>
        );
      }
    },        
    {
      title: '操作',
      dataIndex: 'operation',
      render: (text, record) => {
        return (
          this.state.dataSource.length >= 1
            ? (
              <Popconfirm title='确定删除吗？' onConfirm={() => this.handleDelete(record.sn)}>
                <a href='javascript:;' style={{ color: 'rgb(255, 55, 0)' }}>删除</a>
              </Popconfirm>
            ) : null
        );
      }
    }];

    this.state = {
      dataSource: props.dataSource || [],
      selectedRowKeys: [],
      selectedRows: [],
      loading: false //保存按钮是否激活
    };
  }

  onCheckNeedExIPChange = (v, record) => {
    record.need_extranet_ip = v ? 'yes' : 'no';
    const { dataSource } = this.state;
    dataSource.map(data => {
      if (data.id === record.id) {
        data.need_extranet_ip = v ? 'yes' : 'no';
      }
    });
    this.setState({ dataSource });
  };

  onCheckNeedInIPv6Change = (v, record) => {
    record.need_intranet_ipv6 = v ? 'yes' : 'no';
    const { dataSource } = this.state;
    dataSource.map(data => {
      if (data.id === record.id) {
        data.need_intranet_ipv6 = v ? 'yes' : 'no';
      }
    });
    this.setState({ dataSource });
  };  

  onCheckNeedExIPv6Change = (v, record) => {
    record.need_extranet_ipv6 = v ? 'yes' : 'no';
    const { dataSource } = this.state;
    dataSource.map(data => {
      if (data.id === record.id) {
        data.need_extranet_ipv6 = v ? 'yes' : 'no';
      }
    });
    this.setState({ dataSource });
  };

  handleDelete = (sn) => {
    const dataSource = [...this.state.dataSource];
    this.setState({ dataSource: dataSource.filter(item => item.sn !== sn) });
  };

  execAction = (name, record, selectedRows) => {
    if (name === 'bunchEdit') {
      if (selectedRows.length < 1) {
        return notification.error({ message: '请至少选择一台设备' });
      }
      //清空数据
      this.props.dispatch({
        type: 'bunchEdit/data/clear'
      });
    }

    if (actions[name]) {
      actions[name]({
        record,
        selectedRows,
        dispatch: this.props.dispatch,
        type: this.props.type,
        handleBunchEdit: (onSuccess) => {
          const { sysData, hardwareData, ip, inipv6, exipv6 } = this.props;
          const { selectedRows, dataSource } = this.state;
          //修改硬件配置
          if (hardwareData.id) {
            selectedRows.map(data => {
              data.hardware_template_name = hardwareData.name;
            });
          }
          //修改系统配置
          if (sysData.id) {
            selectedRows.map(data => {
              data.install_type = sysData.install_type;
              data.os_template_name = sysData.name;
            });
          }
          //修改IP
          if (ip) {
            selectedRows.map(data => {
              data.need_extranet_ip = ip;
            });
          }
          if (inipv6) {
            selectedRows.map(data => {
              data.need_intranet_ipv6 = inipv6;
            });
          }
          if (exipv6) {
            selectedRows.map(data => {
              data.need_extranet_ipv6 = exipv6;
            });
          }                    
          //修改dataSource
          let newDataSource = [];
          if (dataSource.length === selectedRows.length) {
            newDataSource = selectedRows;
          } else {
            const copyData = [...dataSource];
            newDataSource = [ ...copyData, ...selectedRows ];
            newDataSource.splice(-1, selectedRows.length);
          }
          this.setState({
            dataSource: newDataSource,
            selectedRowKeys: [],
            selectedRows: []
          });
          onSuccess();
        },
        handleDeviceSubmit: (tableData, onSuccess) => {
          const selectedRows = tableData.selectedRows || [];
          selectedRows.map(item => {
            item.server_room_name = item.server_room.name;
            item.need_extranet_ip = 'no'; //给一个默认值
            item.need_intranet_ipv6 = 'no'; //给一个默认值
            item.need_extranet_ipv6 = 'no'; //给一个默认值
          });
          this.setState({ dataSource: selectedRows });
          onSuccess();
        },
        handleHardwareSubmit: (hardwareData, onSuccess) => {
          const allDevices = this.state.dataSource;
          allDevices.map(data => {
            if (data.id === record.id) {
              data.hardware_template_name = hardwareData.name;
            }
          });
          this.setState({ dataSource: allDevices });
          onSuccess();
        },
        handleSysTemplateSubmit: (sysData, onSuccess) => {
          const allDevices = this.state.dataSource;
          allDevices.map(data => {
            if (data.id === record.id) {
              data.install_type = sysData.install_type;
              data.os_template_name = sysData.name;
            }
          });
          this.setState({ dataSource: allDevices });
          onSuccess();
        }
      });
    }
  };

  getRowSelection = () => {
    const selectedRowKeys = this.state.selectedRowKeys;
    return {
      selectedRowKeys,
      onChange: (selectedRowKeys, selectedRows) => {
        this.setState({
          selectedRowKeys,
          selectedRows
        });
      }
    };
  };

  handleSave = (row) => {
    const newData = [...this.state.dataSource];
    const index = newData.findIndex(item => row.sn === item.sn);
    const item = newData[index];
    newData.splice(index, 1, {
      ...item,
      ...row
    });
    this.setState({ dataSource: newData });
  };

  handleSubmit = () => {
    const newData = [...this.state.dataSource];
    let isReady = false;
    const postData = newData.map(data => {
      if (!data.os_template_name) {
        isReady = false;
        return notification.error({ message: `${data.sn} 请选择操作系统` });
      }
      if (!data.hardware_template_name) {
        isReady = false;
        return notification.error({ message: `${data.sn} 请选择RAID类型` });
      }
      isReady = true;
      return {
        sn: data.sn,
        install_type: data.install_type,
        hardware_template_name: data.hardware_template_name,
        os_template_name: data.os_template_name,
        need_extranet_ip: data.need_extranet_ip,
        need_intranet_ipv6: data.need_intranet_ipv6,
        need_extranet_ipv6: data.need_extranet_ipv6,
      };
    });
    if (isReady) {
      post('/api/cloudboot/v1/devices/settings', postData).then(res => {
        this.setState({ loading: true });
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        //操作成功置灰保存按钮，防止多次点击提交
        //this.setState({ loading: true });
        hashHistory.push('/device/setting');
        // confirm({
        //   title: '操作成功，是否跳转到装机列表页面',
        //   content: '',
        //   okText: '是',
        //   cancelText: '继续操作',
        //   onOk: () => {
        //     hashHistory.push('/device/setting');
        //   },
        //   onCancel: () => {
        //     this.reload();
        //   }
        // });
      });
    }
  };

  reload = () => {
    this.setState({
      dataSource: [],
      selectedRowKeys: [],
      selectedRows: [],
      loading: false
    });
  };


  render() {
    const { dataSource, loading } = this.state;
    const components = {
      body: {
        row: EditableFormRow,
        cell: EditableCell
      }
    };
    const columns = this.columns.map((col) => {
      if (!col.editable) {
        return col;
      }
      return {
        ...col,
        onCell: record => ({
          record,
          editable: col.editable,
          dataIndex: col.dataIndex,
          title: col.title,
          type: col.type,
          placeholder: col.placeholder,
          options: col.options,
          required: col.required,
          addonBefore: col.addonBefore,
          handleSave: this.handleSave
        })
      };
    });
    return (
      <div>
        <div className='operate_btns'>
          {
            this.props.from === 'approval' && this.props.type === 'reinstall' &&
            <Button icon='plus' onClick={() => this.execAction('addDevice')} type='primary' style={{ marginRight: 8 }}>
              选择设备
            </Button>
          }  
          {
            this.props.from === 'approval' && this.props.type === 'recycle' && 
            <Button icon='plus' onClick={() => this.execAction('addRecyclingDevice')} type='primary' style={{ marginRight: 8 }}>
              选择设备[回收中]
            </Button>
          }
        
          {/*<Button icon='plus' onClick={() => addDevice or addRecyclingDevice({*/}
          {/*getServerRoom: true,*/}
          {/*query: {*/}
          {/*operation_status: 'run_with_alarm,run_without_alarm,on_shelve' or 'recycling' */}
          {/*},*/}
          {/*limit: 10,*/}
          {/*handleDeviceSubmit: (tableData, onSuccess) => {*/}
          {/*const selectedRows = tableData.selectedRows || [];*/}
          {/*this.setFormValue(selectedRows);*/}
          {/*onSuccess();*/}
          {/*}*/}
          {/*})} type='primary'*/}
          {/*>*/}
          {/*添加设备*/}
          {/*</Button>*/}

          <Button icon='edit' type='primary' onClick={() => this.execAction('bunchEdit', {}, this.state.selectedRows)} style={{ marginRight: 8 }}>
            批量编辑
          </Button>
        </div>
        <Table
          components={components}
          rowClassName={() => 'editable-row'}
          bordered={true}
          dataSource={dataSource}
          columns={columns}
          pagination={false}
          rowSelection={this.getRowSelection()}
        />
        {
          this.props.from === 'origin' &&
          <div className='pull-right' style={{ marginTop: 16 }}>
            <Button onClick={() => hashHistory.push('/device/list')} style={{ marginRight: 8 }}>
              取消
            </Button>
            <Button onClick={this.handleSubmit} type='primary' disabled={loading}>
              保存
            </Button>
          </div>
        }
        {
          this.props.from === 'approval' && <ApprovalForm type={this.props.type} dataSource={this.state.dataSource} />
        }
      </div>
    );
  }
}

