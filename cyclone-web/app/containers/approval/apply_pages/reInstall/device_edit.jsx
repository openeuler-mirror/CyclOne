import React from 'react';
import { Table, Input, Button, Popconfirm, Form, Modal } from 'antd';
const FormItem = Form.Item;
import actions from 'containers/device/entry/actions';
const EditableContext = React.createContext();
import { notification, Switch } from 'antd';
import { post } from 'common/xFetch2';
const confirm = Modal.confirm;
import { hashHistory } from 'react-router';
import { YES_NO } from "common/enums";
//废弃
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
    this.columns = [
      {
        title: '固资编号',
        dataIndex: 'fixed_asset_number'
      }, {
        title: '序列号',
        dataIndex: 'sn'
      }, {
        title: '操作系统',
        dataIndex: 'os_template_name',
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
            <Button icon='plus' onClick={() => this.execAction('addSystem', record)}>添加操作系统</Button>
        );
        }
      }, {
        title: 'RAID类型',
        dataIndex: 'hardware_template_name',
        required: true,
        render: (text, record) => {
          if (text) {
            return <span className='hover-edit' onClick={() => this.execAction('addHardware', record)}>{text}</span>;
          }
          return (
            <Button icon='plus' onClick={() => this.execAction('addHardware', record)}>添加RAID类型</Button>
        );
        }
      }, {
        title: '是否分配外网IPv4',
        dataIndex: 'need_extranet_ip',
        editable: false,
        required: true,
        type: 'select',
        render: (text, record) => {
          return <Switch checkedChildren='是' unCheckedChildren='否' checked={text === 'yes'} onChange={(value) => this.onCheckChange(value, record)}/>;
        }
      }, {
        title: '是否分配内网IPv6',
        dataIndex: 'need_intranet_ipv6',
        editable: false,
        required: true,
        type: 'select',
        render: (text, record) => {
          return <Switch checkedChildren='是' unCheckedChildren='否' checked={text === 'yes'} onChange={(value) => this.onCheckChange(value, record)}/>;
        }
      }, {
        title: '是否分配外网IPv6',
        dataIndex: 'need_extranet_ipv6',
        editable: false,
        required: true,
        type: 'select',
        render: (text, record) => {
          return <Switch checkedChildren='是' unCheckedChildren='否' checked={text === 'yes'} onChange={(value) => this.onCheckChange(value, record)}/>;
        }
      }, {                
        title: '设备型号',
        dataIndex: 'model'
      }, {
        title: '机位编号',
        dataIndex: 'server_usite',
        render: (text) => <span>{text.number}</span>
      }, {
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

  componentWillReceiveProps(props) {
    if (props.dataSource.length > 0) {
      this.setState({
        dataSource: props.dataSource
      });
    }
  }

  onCheckChange = (v, record) => {
    record.need_extranet_ip = v ? 'yes' : 'no';
    const { dataSource } = this.state;
    dataSource.map(data => {
      if (data.id === record.id) {
        data.need_extranet_ip = v ? 'yes' : 'no';
      }
    });
    this.setState({ dataSource });
    this.props.setFormValue(dataSource);
  };

  handleDelete = (sn) => {
    const dataSource = [...this.state.dataSource];
    const newDataSource = dataSource.filter(item => item.sn !== sn);
    this.setState({ dataSource: newDataSource });
    this.props.setFormValue(newDataSource);
  };

  execAction = (name, record, selectedRows) => {
    if (name === 'bunchEdit') {
      if (selectedRows.length < 1) {
        return notification.error({ message: '请至少选择一台设备' });
      }
      //清空数据
      this.dispatch({
        type: 'bunchEdit/data/clear'
      });
    }

    if (actions[name]) {
      actions[name]({
        record,
        selectedRows,
        dispatch: this.dispatch,
        handleBunchEdit: (onSuccess) => {
          const { sysData, hardwareData, ip } = this.props;
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
        handleHardwareSubmit: (hardwareData, onSuccess) => {
          const allDevices = this.state.dataSource;
          allDevices.map(data => {
            if (data.id === record.id) {
              data.hardware_template_name = hardwareData.name;
            }
          });
          this.setState({ dataSource: allDevices });
          this.props.setFormValue(allDevices);
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
          this.props.setFormValue(allDevices);
          onSuccess();
        }
      });
    }
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
    this.props.setFormValue(newData);
  };

  reload = () => {
    this.setState({
      dataSource: [],
      selectedRowKeys: [],
      selectedRows: [],
      loading: false
    });
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
        <Button icon='edit' type='primary' onClick={() => this.execAction('bunchEdit', {}, this.state.selectedRows)} style={{ marginRight: 8 }}>
          批量编辑
        </Button>
        <Table
          components={components}
          rowClassName={() => 'editable-row'}
          bordered={true}
          dataSource={dataSource}
          columns={columns}
          pagination={false}
          rowSelection={this.getRowSelection()}
        />
      </div>
    );
  }
}

