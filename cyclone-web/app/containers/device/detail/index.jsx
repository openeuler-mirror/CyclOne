import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { Tabs, Icon, Row, Col, Table, Button, Form, Input, Select, notification, DatePicker, Tooltip } from 'antd';
import { renderFormDetail, getBreadcrumb } from 'common/utils';
const TabPane = Tabs.TabPane;
const FormItem = Form.Item;
const Option = Select.Option;
import { PRIVILEGE_LEVEL, TIME_FORMAT, NICSIDE, OPERATION_STATUS, getSearchList, YES_NO, DEVICE_MAINTENANCE_SERVICE_STATUS } from 'common/enums';
import Detail from '../common/detail';
import { put, getWithArgs } from 'common/xFetch2';
import moment from 'moment';
import { getPermissonBtn } from 'common/utils';

class Container extends React.Component {
  state = {
    sn: this.props.params.sn,
    isEdit: false,
    isOobEdit: false,
    isOSSystemInfoEdit: false,
    showPwd: false,
    room: [],
    cabinet: [],
    usite: [],
    store_rooms: [],
    virtualCabinets: []
  };

  componentDidMount() {
    this.reload();
    this.props.dispatch({
      type: 'device-detail/idc/get'
    });
  }

  //编辑时候给机房、机架和机位赋下拉框的值
  UNSAFE_componentWillReceiveProps(props) {
    const { detailInfo } = props.data;
    if (!detailInfo.loading && detailInfo.loading !== this.props.data.detailInfo.loading) {
      const device_page_resp = detailInfo.data.device_page_resp;
      const idc_id = device_page_resp.idc.id;    
      if (idc_id) {
        this.getRoom(idc_id, 1, false);
        const store_room_id = device_page_resp.store_room ? device_page_resp.store_room.id : null;
        if (store_room_id) {
          this.getVirtualCabinets(store_room_id, 1, false);
        }
        const server_room_id = device_page_resp.server_room ? device_page_resp.server_room.id : null;
        if (server_room_id) { 
          this.getCabinet(server_room_id, 1, false);    
          const server_cabinet_id = device_page_resp.server_cabinet ? device_page_resp.server_cabinet.id : null;
          if (server_cabinet_id) {      
            this.getUsite(server_cabinet_id, 1, false);
          }
        }
      }    
    }
  }

  reload = () => {
    this.props.dispatch({
      type: 'device-detail/detail-info/get',
      payload: this.state.sn
    });
  };

  //保存修改信息
  handleSubmit = () => {
    const { data } = this.props.data.detailInfo;
    this.props.form.validateFields((err, values) => {
      if (err) {
        return notification.error({ message: '还有未填项' });
      }
      delete values.created_at;
      delete values.updated_at;
      delete values.oob_password;
      delete values.oob_user;
      const postValues = {
        ...values,
        id: data.device_page_resp.id,
        'started_at': values['started_at'] ? values['started_at'].format(TIME_FORMAT) : '',
        'onshelve_at': values['onshelve_at'] ? values['onshelve_at'].format(TIME_FORMAT) : ''
      };
      put(`/api/cloudboot/v1/device`, postValues).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        notification.success({ message: res.message });
        this.setState({ isEdit: false });
        this.reload();
      });
    });
  };
  reloadOob = () => {
    const { data } = this.props.data.detailInfo;
    const { setFieldsValue } = this.props.form;
    setFieldsValue({ oob_user: data.device_page_resp.oob_user, oob_password: data.device_page_resp.oob_password });
  };
  //保存带外信息
  handleOobSubmit = () => {
    this.props.form.validateFields((err, values) => {
      const postValues = {
        sn: this.state.sn,
        oob_user_name: values.oob_user,
        oob_password_old: values.oob_password
      };
      put(`/api/cloudboot/v1/devices/${this.state.sn}/oob/password`, postValues).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        notification.success({ message: res.message });
        this.setState({ isOobEdit: false });
        this.reload();
      });
    });
  };

  reloadOS = () => {
    const { data } = this.props.data.detailInfo;
    const { setFieldsValue } = this.props.form;
    setFieldsValue({ os: data.os});
  };
  //保存带外信息
  handleOSSystemInfoSubmit = () => {
    this.props.form.validateFields((err, values) => {
      const postValues = {
        sn: this.state.sn,
        os: values.os,
      };
      put(`/api/cloudboot/v1/devices/settings`, postValues).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        notification.success({ message: res.message });
        this.setState({ isOSSystemInfoEdit: false });
        this.reload();
      });
    });
  };  

  clearStoreRoom = () => {
    const { setFieldsValue, getFieldValue } = this.props.form;
    const room = getFieldValue('store_room_id');
    if (room) {
      setFieldsValue({ store_room_id: null });
    }
    this.setState({
      store_rooms: []
    });
  };

  clearVirtualCabinet = () => {
    const { setFieldsValue, getFieldValue } = this.props.form;
    const room = getFieldValue('virtual_cabinet_id');
    if (room) {
      setFieldsValue({ virtual_cabinet_id: null });
    }
    this.setState({
      virtualCabinets: []
    });
  };


  clearRoom = () => {
    const { setFieldsValue, getFieldValue } = this.props.form;
    const room = getFieldValue('server_room_id');
    if (room) {
      setFieldsValue({ server_room_id: null });
    }
    this.setState({
      room: []
    });
  };
  clearCabinet = () => {
    const { setFieldsValue, getFieldValue } = this.props.form;
    const cabinet = getFieldValue('server_cabinet_id');
    if (cabinet) {
      setFieldsValue({ server_cabinet_id: null });
    }
    this.setState({
      cabinet: []
    });
  };
  clearUsite = () => {
    const { setFieldsValue, getFieldValue } = this.props.form;
    const usite = getFieldValue('server_usite_id');
    if (usite) {
      setFieldsValue({ server_usite_id: null });
    }
    this.setState({
      usite: []
    });
  };

  //选择idc获取机房列表
  getRoom = (v, page, clear) => {
    this.getRoomPage(v, page);
    this.getStoreRoomPage(v, page);
    if (clear) {
      this.clearRoom();
      this.clearCabinet();
      this.clearUsite();
      this.clearStoreRoom();
      this.clearVirtualCabinet();
    }
  };

  getStoreRoomPage = (v, page) => {
    getWithArgs('/api/cloudboot/v1/store-rooms', { page: page, page_size: 100, idc_id: v }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          store_rooms: [ ...preState.store_rooms, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getStoreRoomPage(v, page + 1);
        }
      });
    });
  };

  getRoomPage = (v, page) => {
    getWithArgs('/api/cloudboot/v1/server-rooms', { page: page, page_size: 100, idc_id: v }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          room: [ ...preState.room, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getRoomPage(v, page + 1);
        }
      });
    });
  };

  //选择机房获取机架列表
  getCabinet = (v, page, clear) => {
    this.getCabinetPage(v, page);
    if (clear) {
      this.clearCabinet();
      this.clearUsite();
    }
  };

  getCabinetPage = (v, page) => {
    getWithArgs('/api/cloudboot/v1/server-cabinets', { page: page, page_size: 1000, server_room_id: v }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          cabinet: [ ...preState.cabinet, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getCabinetPage(v, page + 1);
        }
      });
    });
  };

    //选择库房获取虚机货架
  getVirtualCabinets = (v, page, clear) => {
    this.getVirtualCabinetsPage(v, page);
    if (clear) {
      this.clearVirtualCabinet();
    }
  };
  getVirtualCabinetsPage = (v, page) => {
    getWithArgs('/api/cloudboot/v1/virtual-cabinets', { page: page, page_size: 100, store_room_id: v }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          virtualCabinets: [ ...preState.virtualCabinets, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getVirtualCabinetsPage(v, page + 1);
        }
      });     
    });
  };
  

  //选择机架获取机位列表
  getUsite = (v, page, clear) => {
    this.getUsitePage(v, page);
    if (clear) {
      this.clearUsite();
    }
  };

  getUsitePage = (v, page) => {
    getWithArgs('/api/cloudboot/v1/server-usites', { page: page, page_size: 100, server_cabinet_id: v, status: 'free' }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      this.setState((preState) => {
        return {
          usite: [ ...preState.usite, ...res.content.records ]
        };
      }, () => {
        if (res.content.total_pages > page) {
          this.getUsitePage(v, page + 1);
        }
      });
    });
  };

  render() {
    const { detailInfo, idc } = this.props.data;
    //console.log("show detailInfo 2:",detailInfo);
    const { data } = detailInfo;
    const { device_page_resp, cpu, memory, nic, oob, pci, hba, disk_slot, raid, motherboard, device_lifecycle_detail_page } = data;
    const { isEdit, isOobEdit, isOSSystemInfoEdit, showPwd } = this.state;
    const { getFieldDecorator } = this.props.form;
    const isSpecial = device_page_resp.usage === '特殊设备';

    return (
      <Layout>
        <Tabs type='card' tabBarExtraContent={getBreadcrumb(this.state.sn)}>
          <TabPane tab='基本信息' key='info'>
            <div className='pull-right'>
              {isEdit ?
                <div>
                  <Button onClick={() => this.setState({ isEdit: false })} style={{ marginRight: 8 }}>取消</Button>
                  <Button onClick={() => this.handleSubmit()} type='primary' >保存</Button>
                </div> :
                <Button disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_update')}
                  onClick={() => this.setState({ isEdit: true })}
                ><Icon type='edit' theme='outlined' />修改设备信息</Button>
              }
            </div>
            <h3 className='detail-title'>设备信息</h3>
            <div className='detail-info'>
              {isEdit ?
                <div>
                  <Form>
                    <Row gutter={24}>
                      <Col span={8}>
                        <FormItem label={'固资编号'}>
                          {getFieldDecorator('fixed_asset_number', {
                            initialValue: device_page_resp.fixed_asset_number,
                            rules: [
                              {
                                required: true,
                                message: '请填写固资编号'
                              }
                            ]
                          })(
                            <Input />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'序列号'}>
                          {getFieldDecorator('sn', {
                            initialValue: device_page_resp.sn,
                            rules: [
                              {
                                required: true,
                                message: '请填写序列号'
                              }
                            ]
                          })(
                            <Input />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'设备型号'}>
                          {getFieldDecorator('model', {
                            initialValue: device_page_resp.model,
                            rules: [
                              {
                                required: true,
                                message: '请填写设备型号'
                              }
                            ]
                          })(
                            <Input />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'设备类型'}>
                          {getFieldDecorator('category', {
                            initialValue: device_page_resp.category,
                            rules: [
                              {
                                required: true,
                                message: '请填写设备类型'
                              }
                            ]
                          })(
                            <Input disabled={isSpecial} />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'厂商'}>
                          {getFieldDecorator('vendor', {
                            initialValue: device_page_resp.vendor,
                            rules: [
                              {
                                required: true,
                                message: '请填写厂商'
                              }
                            ]
                          })(
                            <Input />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'用途'}>
                          {getFieldDecorator('usage', {
                            initialValue: device_page_resp.usage,
                            rules: [
                              {
                                required: true,
                                message: '请填写用途'
                              }
                            ]
                          })(
                            <Input  disabled={isSpecial}/> 
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'运营状态'}>
                          {getFieldDecorator('operation_status', {
                            initialValue: device_page_resp.operation_status,
                            rules: [
                              {
                                required: true,
                                message: '请选择运营状态'
                              }
                            ]
                          })(
                            <Select>
                              {
                                getSearchList(OPERATION_STATUS).map((it) => <Option value={it.value} key={it.value}>{it.label}</Option>)
                              }
                            </Select>
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label='数据中心' >
                          {getFieldDecorator('idc_id', {
                            initialValue: device_page_resp.idc.id || null,
                            rules: [
                              {
                                required: true,
                                message: '请选择数据中心'
                              }
                            ]
                          })(
                            <Select
                              showSearch={true}
                              filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                              onChange={(value) => this.getRoom(value, 1, true)}
                            >
                              {
                                !idc.loading &&
                                idc.data.map((os, index) => <Option value={os.id} key={os.id}>{os.name}</Option>)
                              }
                            </Select>
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'机房管理单元'}>
                          {getFieldDecorator('server_room_id', {
                            initialValue: device_page_resp.server_room.id || null
                            // rules: [
                            //   {
                            //     required: true,
                            //     message: '请选择机房管理单元'
                            //   }
                            // ]
                          })(
                            (<Select
                              showSearch={true}
                              filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                              onChange={(value) => this.getCabinet(value, 1, true)}
                            >
                              {
                                this.state.room.map((os, index) => <Option value={os.id} key={os.id}>{os.name}</Option>)
                              }
                            </Select>)
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label='机架编号' >
                          {getFieldDecorator('server_cabinet_id', {
                            initialValue: device_page_resp.server_cabinet.id || null
                            // rules: [
                            //   {
                            //     required: true,
                            //     message: '请选择机架编号'
                            //   }
                            // ]
                          })(<Select
                            showSearch={true}
                            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                            onChange={(value) => this.getUsite(value, 1, true)}
                          >
                            {
                              this.state.cabinet.map((os, index) => <Option value={os.id} key={os.id}>{os.number}</Option>)
                            }
                          </Select>)}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label='机位编号' >
                          {getFieldDecorator('server_usite_id', {
                            initialValue: device_page_resp.server_usite ? device_page_resp.server_usite.id : null
                            // rules: [
                            //   {
                            //     required: true,
                            //     message: '请选择机位编号'
                            //   }
                            // ]
                          })(<Select
                            showSearch={true}
                            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                          >
                            {
                              this.state.usite.map((os, index) => <Option value={os.id} key={os.id}>{os.number}</Option>)
                            }
                          </Select>)}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label='库房管理单元' >
                          {getFieldDecorator('store_room_id', {
                            initialValue: device_page_resp.store_room ? device_page_resp.store_room.id : null
                            // rules: [
                            //   {
                            //     required: true,
                            //     message: '请选择机位编号'
                            //   }
                            // ]
                          })(<Select
                            onChange={(value) => this.getVirtualCabinets(value, 1, true)}
                            showSearch={true}
                            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                          >
                            {
                              (this.state.store_rooms || []).map((os, index) => <Option value={os.id} key={os.id}>{os.name}</Option>)
                            }
                          </Select>)}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label='虚拟货架' >
                          {getFieldDecorator('virtual_cabinet_id', {
                            initialValue: device_page_resp.virtual_cabinets ? device_page_resp.virtual_cabinets.id : null
                            // rules: [
                            //   {
                            //     required: true,
                            //     message: '请选择机位编号'
                            //   }
                            // ]
                          })(<Select
                            showSearch={true}
                            filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
                          >
                            {
                              (this.state.virtualCabinets || []).map((os, index) => <Option value={os.id} key={os.id}>{os.number}</Option>)
                            }
                          </Select>)}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'代理节点'}>
                          {getFieldDecorator('origin_node', {
                            initialValue: device_page_resp.origin_node
                          })(
                            <Input disabled={true} />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'代理节点IP'}>
                          {getFieldDecorator('origin_node_ip', {
                            initialValue: device_page_resp.origin_node_ip
                          })(
                            <Input disabled={true} />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'启用时间'}>
                          {getFieldDecorator('started_at', {
                            initialValue: device_page_resp.started_at ? moment(device_page_resp.started_at) : null
                          })(
                            <DatePicker showTime={true} format={TIME_FORMAT} style={{ width: '100%' }} />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'上架时间'}>
                          {getFieldDecorator('onshelve_at', {
                            initialValue: device_page_resp.onshelve_at ? moment(device_page_resp.onshelve_at) : null
                          })(
                            <DatePicker showTime={true} format={TIME_FORMAT} style={{ width: '100%' }} />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'创建时间'}>
                          {getFieldDecorator('created_at', {
                            initialValue: device_page_resp.created_at
                          })(
                            <Input disabled={true} />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'更新时间'}>
                          {getFieldDecorator('updated_at', {
                            initialValue: device_page_resp.updated_at
                          })(
                            <Input disabled={true} />
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'RAID结构'}>
                          {getFieldDecorator('raid_remark', {
                            initialValue: device_page_resp.raid_remark
                          })(
                            <Input/>
                          )}
                        </FormItem>
                      </Col>
                      <Col span={8}>
                        <FormItem label={'硬件说明'}>
                          {getFieldDecorator('hardware_remark', {
                            initialValue: device_page_resp.hardware_remark
                          })(
                            <Input.TextArea />
                          )}
                        </FormItem>
                      </Col>
                    </Row>
                  </Form>
                </div> :
                <Row>
                  <Detail sn={this.state.sn} data={detailInfo.data} />
                </Row>
              }
            </div>
            <h3 className='detail-title'>维保信息</h3>
            <div className='detail-info'>
              <Row>
               {renderFormDetail([
                  {
                    label: '资产归属',
                    value: device_lifecycle_detail_page.asset_belongs
                  },
                  {
                    label: '负责人',
                    value: device_lifecycle_detail_page.owner
                  },          
                  {
                    label: '是否租赁',
                    value: YES_NO[device_lifecycle_detail_page.is_rental]
                  },
                  {
                    label: '维保服务供应商',
                    value: device_lifecycle_detail_page.maintenance_service_provider
                  },
                  {
                    label: '维保服务起始日期',
                    value: device_lifecycle_detail_page.maintenance_service_date_begin
                  },
                  {
                    label: '维保服务截止日期',
                    value: device_lifecycle_detail_page.maintenance_service_date_end
                  },
                  {
                    label: '维保服务状态',
                    value: DEVICE_MAINTENANCE_SERVICE_STATUS[device_lifecycle_detail_page.maintenance_service_status]
                  },
                  {
                    label: '退役日期',
                    value: device_lifecycle_detail_page.device_retired_date
                  }, 
                ])}
                </Row>
                <Row>
                <Col span={8}>
                  <FormItem label={'维保服务内容'}>
                    {getFieldDecorator('maintenance_service', {
                      initialValue: device_lifecycle_detail_page.maintenance_service
                    })(
                      <Input.TextArea disabled={true} autoSize={ true } />
                    )}
                  </FormItem>
                </Col>
                <Col span={8}>  
                  <FormItem label={'物流服务内容'}>
                    {getFieldDecorator('logistics_service', {
                      initialValue: device_lifecycle_detail_page.logistics_service
                    })(
                      <Input.TextArea disabled={true} autoSize={ true } />
                    )}
                  </FormItem>
                </Col>
              </Row>
            </div>
            <div className='pull-right'>
              {isOobEdit ?
                <div>
                  <Button onClick={() => {
                    this.setState({ isOobEdit: false });
                    this.reloadOob();
                  }} style={{ marginRight: 8 }}
                  >取消</Button>
                  <Button onClick={() => this.handleOobSubmit()} type='primary' >保存</Button>
                </div> :
                <Button disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_update_oob')}
                  onClick={() => this.setState({ isOobEdit: true })}
                ><Icon type='edit' theme='outlined' />修改带外信息</Button>
              }
            </div>
            <h3 className='detail-title'>带外信息</h3>
            <div className='detail-info'>
              <Row>
                <Col span={8}>
                  <FormItem label={'带外用户名'}>
                    {getFieldDecorator('oob_user', {
                      initialValue: device_page_resp.oob_user
                    })(
                      <Input disabled={!isOobEdit} />
                    )}
                  </FormItem>
                </Col>
                <Col span={8}>
                  <FormItem label={isOobEdit ? '带外密码（旧）' : '带外密码'}>
                    {getFieldDecorator('oob_password', {
                      initialValue: device_page_resp.oob_password
                    })(
                      <Input disabled={!isOobEdit} type={showPwd ? 'text' : 'password'} addonAfter={<Icon type='eye' onClick={() => this.setState({ showPwd: !this.state.showPwd })} />} />
                    )}
                  </FormItem>
                </Col>
                {
                  isOobEdit &&
                  <Col span={8}>
                    <FormItem label={'带外密码（新）'}>
                      <span>新密码默认随机生成</span>
                    </FormItem>
                  </Col>
                }

              </Row>
            </div>
            <div className='pull-right'>
              {isOSSystemInfoEdit ?
                <div>
                  <Button onClick={() => {
                    this.setState({ isOSSystemInfoEdit: false });
                    this.reloadOS();
                  }} style={{ marginRight: 8 }}
                  >取消</Button>
                  <Button onClick={() => this.handleOSSystemInfoSubmit()} type='primary' >保存</Button>
                </div> :
                <Button disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_physical_machine_update_oob')}
                  onClick={() => this.setState({ isOSSystemInfoEdit: true })}
                ><Icon type='edit' theme='outlined' />修改系统信息</Button>
              }
            </div>
            <h3 className='detail-title'>系统信息</h3>
            <div className='detail-info'>
              <Row>
                {renderFormDetail([
                  {
                    label: '内网 IP',
                    value: device_page_resp.intranet_ip
                  },
                  {
                    label: '外网 IP',
                    value: device_page_resp.extranet_ip
                  }
                ])}
              </Row>
              <Row>
              {renderFormDetail([
                  {
                    label: '内网 IPv6',
                    value: device_page_resp.intranet_ipv6
                  },
                  {
                    label: '外网 IPv6',
                    value: device_page_resp.extranet_ipv6
                  }                
                  //{
                  //  label: '操作系统',
                  //  value: data.os
                  //}
                ])}
              </Row>
              <Row>
                <Col span={8}>
                  <FormItem label={isOSSystemInfoEdit ? '操作系统（旧）' : '操作系统'}>
                    {getFieldDecorator('os', {
                      initialValue: data.os
                    })(
                      <Input disabled={!isOSSystemInfoEdit}/>
                    )}
                  </FormItem>
                </Col>
                {
                  isOSSystemInfoEdit &&
                  <Col span={8}>
                    <FormItem label={'操作系统（新）'}>
                      <span>根据新操作系统名称关联PXE装机配置模板（不存在则新增）</span>
                    </FormItem>
                  </Col>
                }
              </Row>
            </div>
          </TabPane>
          <TabPane tab='CPU' key='cpu'>
            <Table columns={this.getCpuColumns()} dataSource={cpu.physicals} pagination={false} footer={() => `总计核数： ${cpu.total_cores}`} />
          </TabPane>
          <TabPane tab='内存' key='memory'>
            <Table columns={this.getMemoryColumns()} dataSource={memory.items} pagination={false} footer={() => `总计容量： ${memory.total_size_mb} MB`} />
          </TabPane>
          <TabPane tab='存储' key='disk'>
            {raid.items && raid.items.length > 0 ?
              <Table defaultExpandAllRows={true} columns={this.getRaidColumns()} title={() => `RAID卡`} expandedRowRender={this.expandedRowRender} dataSource={raid.items} pagination={false} />
              :
              <Table title={() => `磁盘`} columns={this.getDiskColumns()} dataSource={disk_slot.items} pagination={false} />
            }
          </TabPane>
          <TabPane tab='网卡' key='nic'>
            <Table columns={this.getNicColumns()} dataSource={nic.items} pagination={false} />
          </TabPane>
          <TabPane tab='主板' key='motherboard'>
            <Table columns={this.getMotherboardColumns()} dataSource={[motherboard]} pagination={false} />
          </TabPane>
          <TabPane tab='带外' key='oob'>
            <Table columns={this.getOobNetworkColumns()} dataSource={oob.network ? [oob.network] : []} pagination={false} title={() => '网络'} />
            <Table columns={this.getOobUserColumns()} dataSource={oob.user || []} pagination={false} title={() => '用户'} />
            <div className='ant-table-title'>
              固件版本
            </div>
            <div>{oob.firmware}</div>
          </TabPane>
          <TabPane tab='PCI' key='pci'>
            <Table columns={this.getPciColumns()} dataSource={pci.slots} pagination={false} footer={() => `插槽总数： ${pci.total_slots} 个`} />
          </TabPane>
          <TabPane tab='HBA' key='hba'>
            <Table columns={this.getHbaColumns()} dataSource={hba.items} pagination={false} />
          </TabPane>
          <TabPane tab='生命周期记录' key='lifecycle_log'>
            <Table columns={this.getLifecycleLogColumns()} dataSource={device_lifecycle_detail_page.lifecycle_log} pagination={false}  />
          </TabPane>
        </Tabs>
      </Layout>
    );
  }

  getCpuColumns = () => {
    return [
      {
        title: 'CPU型号',
        dataIndex: 'model_name',
        key: 'model_name',
        width: '50%'
      },
      {
        title: 'CPU核数',
        dataIndex: 'cores',
        key: 'cores',
        width: '50%'
      }
    ];
  };
  getDiskColumns = () => {
    return [
      {
        title: '物理驱动器',
        dataIndex: 'name',
        key: 'name',
        width: '30%'
      },
      {
        title: '容量',
        dataIndex: 'raw_size',
        key: 'raw_size',
        width: '30%'
      },
      {
        title: '类型',
        dataIndex: 'media_type',
        key: 'media_type',
        width: '40%'
      }
    ];
  };
  getRaidColumns = () => {
    return [
      {
        title: 'ID',
        dataIndex: 'id',
        key: 'id',
        width: '30%'
      },
      {
        title: '型号',
        dataIndex: 'model_name',
        key: 'model_name',
        width: '30%'
      },
      {
        title: '固件版本',
        dataIndex: 'firmware_version',
        key: 'firmware_version',
        width: '40%'
      }
    ];
  };
  getMemoryColumns = () => {
    return [
      {
        title: '内存卡',
        dataIndex: 'locator',
        key: 'locator',
        width: '25%'
      },
      {
        title: '容量',
        dataIndex: 'size',
        key: 'size',
        width: '25%'
      },
      {
        title: '速率',
        dataIndex: 'speed',
        key: 'speed',
        width: '25%'
      },
      {
        title: '类型',
        dataIndex: 'type',
        key: 'type',
        width: '25%'
      }
    ];
  };
  getLifecycleLogColumns = () => {
    return [
      {
        title: '操作类型',
        dataIndex: 'operation_type',
        key: 'operation_type',
        width: '10%'
      },
      {
        title: '操作内容',
        dataIndex: 'operation_detail',
        key: 'operation_detail',
        width: '60%',
        render: (text) => <Tooltip placement="leftTop" title={text}>{text}</Tooltip>
      },
      {
        title: '操作用户',
        dataIndex: 'operation_user',
        key: 'operation_user',
        width: '10%'
      },
      {
        title: '操作时间',
        dataIndex: 'operation_time',
        key: 'operation_time',
        width: '20%'
      }
    ];
  };  
  getNicColumns = () => {
    return [
      {
        title: '名称',
        dataIndex: 'name',
        key: 'name'
      },
      {
        title: 'Mac',
        dataIndex: 'mac',
        key: 'mac'
      },
      {
        title: 'BootOS IP',
        dataIndex: 'ip',
        key: 'ip'
      },
      {
        title: '设备总线',
        dataIndex: 'businfo',
        key: 'businfo'
      },
      {
        title: '槽位',
        dataIndex: 'designation',
        key: 'designation'
      },
      {
        title: '内/外置',
        dataIndex: 'side',
        key: 'side',
        render: (text) => {
          return NICSIDE[text];
        }
      },
      {
        title: '速率',
        dataIndex: 'speed',
        key: 'speed'
      },
      {
        title: '厂商',
        dataIndex: 'company',
        key: 'company'
      },
      {
        title: '型号',
        dataIndex: 'model_name',
        key: 'model_name'
      },
      {
        title: '固件版本',
        dataIndex: 'firmware_version',
        key: 'firmware_version'
      }
    ];
  };
  getPciColumns = () => {
    return [
      {
        title: '槽位',
        dataIndex: 'designation',
        key: 'designation',
        width: '30%'
      },
      {
        title: '设备类型',
        dataIndex: 'type',
        key: 'type',
        width: '30%'
      },
      {
        title: '当前使用情况',
        dataIndex: 'current_usage',
        key: 'current_usage',
        width: '40%'
      }
    ];
  };
  getHbaColumns = () => {
    return [
      {
        title: '主机',
        dataIndex: 'host',
        key: 'host',
        width: '30%'
      },
      {
        title: 'WWPN',
        dataIndex: 'wwpns',
        key: 'wwpns',
        width: '30%'
      },
      {
        title: 'WWNN',
        dataIndex: 'wwnns',
        key: 'wwnns',
        width: '40%'
      }
    ];
  };
  getOobNetworkColumns = () => {
    return [
      {
        title: 'IP 来源',
        dataIndex: 'ip_src',
        key: 'ip_src',
        width: '25%'
      },
      {
        title: 'IP',
        dataIndex: 'ip',
        key: 'ip',
        width: '25%'
      },
      {
        title: '掩码',
        dataIndex: 'netmask',
        key: 'netmask',
        width: '25%'
      },
      {
        title: '网关',
        dataIndex: 'gateway',
        key: 'gateway',
        width: '25%'
      }
    ];
  };
  getOobUserColumns = () => {
    return [
      {
        title: 'ID',
        dataIndex: 'id',
        key: 'id',
        width: '30%'
      },
      {
        title: '用户名',
        dataIndex: 'name',
        key: 'name',
        width: '30%'
      },
      {
        title: '权限级别',
        dataIndex: 'privilege_level',
        key: 'privilege_level',
        width: '40%',
        render: (text) => {
          return PRIVILEGE_LEVEL[text];
        }
      }
    ];
  };
  getMotherboardColumns = () => {
    return [
      {
        title: '厂商',
        dataIndex: 'manufacturer',
        key: 'manufacturer',
        width: '30%'
      },
      {
        title: '产品型号',
        dataIndex: 'product_name',
        key: 'product_name',
        width: '30%'
      },
      {
        title: '序列号',
        dataIndex: 'serial_number',
        key: 'serial_number',
        width: '40%'
      }
    ];
  };
  //Raid卡内置网卡的表格渲染
  expandedRowRender = (record) => {
    const { detailInfo } = this.props.data;
    const { disk_slot } = detailInfo.data;
    const diskList = [];
    if (disk_slot.items) {
      disk_slot.items.forEach(data => {
        if (data.controller_id === record.id) {
          diskList.push(data);
        }
      });
    }
    return <Table title={() => `磁盘`} columns={this.getDiskColumns()} dataSource={diskList} pagination={false} />;
  };
}
function mapStateToProps(state) {
  return {
    data: state.get('device-detail').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(Container));
