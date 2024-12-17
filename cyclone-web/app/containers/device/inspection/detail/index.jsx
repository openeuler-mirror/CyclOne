import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { Tabs, Table, Form, Select, Tooltip} from 'antd';
import { renderFormDetail, getBreadcrumb, geTabsTitle } from 'common/utils';
const TabPane = Tabs.TabPane;
const Option = Select.Option;
import { put } from 'common/xFetch2';
import { PRIVILEGE_LEVEL, IP_REGEXP, NICSIDE } from 'common/enums';
import Detail from 'containers/device/common/detail';
import { INSPECTION_RESULT } from 'common/enums';
import { hashHistory } from 'react-router';


class Container extends React.Component {
  state = {
    sn: this.props.params.sn,
    startTime: ''
  };

  componentDidMount() {
    this.reload();
  }

  reload = () => {
    this.props.dispatch({
      type: 'inspection-detail/detail-info/get',
      payload: { sn: this.state.sn }
    });
    this.props.dispatch({
      type: 'inspection-detail/start-time/get',
      payload: this.state.sn
    });
  };

  getResult = (v) => {
    this.props.dispatch({
      type: 'inspection-detail/detail-info/get',
      payload: {
        sn: this.state.sn,
        start_time: v
      }
    });
    this.setState({ startTime: v });
  };

  getColumns = (type) => {
    if (type == "sensor") {
      return [
        {
          title: '传感器',
          dataIndex: 'name',
          key: 'name',
          width: '20%'
        },
        {
          title: '类型',
          dataIndex: 'type',
          key: 'type',
          width: '20%'
        },
        {
          title: '状态',
          dataIndex: 'state',
          key: 'state',
          width: '20%',
          render: (text, record) => {
            if (text === 'Critical') {
              return <span style={{ color: '#ff3700' }}>{text}</span>;
            }
            return text;
          }
        },
        {
          title: '读数',
          dataIndex: 'reading',
          key: 'reading',
          width: '20%',
          render: (text, record) => {
            return <span>{text} {record.units}</span>;
          }
        },
        {
          title: '事件',
          dataIndex: 'event',
          key: 'event',
          width: '20%',
          render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
        }
      ];
    } else if (type == "sel"){
      return [
        {
          title: '日期',
          dataIndex: 'date',
          key: 'date',
          width: '10%',
        },
        {
          title: '时间',
          dataIndex: 'time',
          key: 'time',
          width: '10%'
        },
        {
          title: '类型',
          dataIndex: 'type',
          key: 'type',
          width: '20%'
        },
        {
          title: '描述',
          dataIndex: 'name',
          key: 'name',
          width: '20%'
        },
        {
          title: '事件',
          dataIndex: 'event',
          key: 'event',
          width: '20%',
          render: (text) => <Tooltip placement="top" title={text}>{text}</Tooltip>
        },
        {
          title: '状态',
          dataIndex: 'state',
          key: 'state',
          width: '20%',
          render: (text, record) => {
            if (text === 'Critical') {
              return <span style={{ color: '#ff3700' }}>{text}</span>;
            }
            return text;
          }
        }
      ];
    }

  };

  renderContent = (loading, dataSource) => {
    const { startTime } = this.props;
    return (
      <div>
        <div className='buttons-margin'>
          巡检时间：
          <Select style={{ width: 200 }} onChange={this.getResult} value={this.state.startTime || startTime.data[0]} >
            {
              startTime.data.map(item => <Option value={item} key={item}>{item}</Option>)
            }
          </Select>
        </div>
        <Table columns={this.getColumns("sensor")} loading={loading} dataSource={dataSource} pagination={{ showSizeChanger: true, showQuickJumper: true }} />
      </div>
    );
  };

  renderSelContent = (loading, dataSource) => {
    const { startTime } = this.props;
    return (
      <div>
        <div className='buttons-margin'>
          巡检时间：
          <Select style={{ width: 200 }} onChange={this.getResult} value={this.state.startTime || startTime.data[0]} >
            {
              startTime.data.map(item => <Option value={item} key={item}>{item}</Option>)
            }
          </Select>
        </div>
        <Table columns={this.getColumns("sel")} loading={loading} dataSource={dataSource} pagination={{ showSizeChanger: true, showQuickJumper: true }} />
      </div>
    );
  };

  renderTab = (num, title) => {
    if (num > 0) {
      return <span>{title}<span className='spanCircle fill red'>{num}</span></span>;
    }
    return title;
  };
  onTabChange = (key) => {
    const currentQuery = this.props.location.query;
    currentQuery.type = key;
    hashHistory.push(this.props.location);
  };
  render() {
    const { detailInfo } = this.props;
    const { loading } = detailInfo;
    const temperature = [], voltage = [], fan = [], Memory = [], power_supply = [], others = [], sel = [];
    let temperature_error = 0, voltage_error = 0, fan_error = 0, Memory_error = 0, power_supply_error = 0, others_error = 0, sel_error = 0;
    const inspectData = detailInfo.data.result || [];
    const ipmiSelData = detailInfo.data.selresult || [];
    inspectData.forEach(item => {
      if (item.type === 'Temperature') {
        temperature.push(item);
        if (item.state === 'Critical') {
          temperature_error ++;
        }
      } else if (item.type === 'Fan') {
        fan.push(item);
        if (item.state === 'Critical') {
          fan_error ++;
        }
      } else if (item.type === 'Cooling Device') {
          fan.push(item);
          if (item.state === 'Critical') {
            fan_error ++;
          }        
      } else if (item.type === 'Power Supply') {
        power_supply.push(item);
        if (item.state === 'Critical') {
          power_supply_error ++;
        }
      } else if (item.type === 'Memory') {
        Memory.push(item);
        if (item.state === 'Critical') {
          Memory_error ++;
        }
      } else if (item.type === 'Voltage') {
        voltage.push(item);
        if (item.state === 'Critical') {
          voltage_error ++;
        }
      } else if (item.type === 'Voltage') {
        voltage.push(item);
        if (item.state === 'Critical') {
          voltage_error ++;
        }
      } else {
        others.push(item);
        if (item.state === 'Critical') {
          others_error ++;
        }
      }
    });

    ipmiSelData.forEach(item => {
      sel.push(item);
      if (item.state === 'Critical') {
        sel_error ++;
        } 
    });

    const { query } = this.props.location;
    const defaultKey = query.type ? query.type : 'temperature';
    return (
      <Layout>
        {getBreadcrumb(this.state.sn, '/device/inspection/list')}
        <div className='detail-body' style={{ marginTop: 0 }}>
          <h3 className='detail-title'>设备信息</h3>
          <div className='detail-info'>
            <Detail sn={this.state.sn} />
          </div>
          <Tabs type='card' onChange={this.onTabChange} activeKey={defaultKey} defaultActiveKey={defaultKey} tabBarExtraContent={geTabsTitle('巡检详情[No Nominal]')}>
            <TabPane tab={this.renderTab(temperature_error, '温度')} key='temperature'>
              {this.renderContent(loading, temperature)}
            </TabPane>
            <TabPane tab={this.renderTab(voltage_error, '电压')} key='voltage'>
              {this.renderContent(loading, voltage)}
            </TabPane>
            <TabPane tab={this.renderTab(fan_error, '风扇')} key='fan'>
              {this.renderContent(loading, fan)}
            </TabPane>
            <TabPane tab={this.renderTab(Memory_error, '内存')} key='memory'>
              {this.renderContent(loading, Memory)}
            </TabPane>
            <TabPane tab={this.renderTab(power_supply_error, '电源')} key='power_supply'>
              {this.renderContent(loading, power_supply)}
            </TabPane>
            <TabPane tab={this.renderTab(sel_error, '事件')} key='sel'>
              {this.renderSelContent(loading, sel)}
            </TabPane>
            <TabPane tab={this.renderTab(others_error, '其它')} key='others'>
              {this.renderContent(loading, others)}
            </TabPane>
          </Tabs>
        </div>
      </Layout>
    );
  }
}
function mapStateToProps(state) {
  return {
    detailInfo: state.get('device-inspection-detail').toJS().detailInfo,
    //deviceInfo: state.get('device-inspection-detail').toJS().deviceInfo,
    startTime: state.get('device-inspection-detail').toJS().startTime,
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Form.create()(Container));
