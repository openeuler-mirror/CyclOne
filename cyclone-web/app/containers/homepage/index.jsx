import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Col, Row, Card, Radio } from 'antd';
import Line from './components/line';
import Pine from './components/pie';
const RadioButton = Radio.Button;
const RadioGroup = Radio.Group;
import { Link } from 'react-router';
import { arrayFind } from 'common/util';

class Container extends Component {
  constructor(props) {
    super(props);
  }
  componentDidMount() {
    this.props.dispatch({
      type: 'homepage/devices/get'
    });
    this.props.dispatch({
      type: 'homepage/inspections/get',
      payload: {
        period: 'latest_week'
      }
    });
  }


  onChange = (e) => {
    this.props.dispatch({
      type: 'homepage/inspections/get',
      payload: {
        period: e.target.value
      }
    });
  };

  render() {
    const { devices, inspections } = this.props.data;
    const { permissions } = this.props.userInfo;
    let ifPermise;
    if (permissions) {
      ifPermise = arrayFind(permissions, 'menu_device_setting');
    }
    return (
      <div className='home-page'>
        <div className='home-main'>
          <Row gutter={16} className='margin-bottom'>
            <Col span='6'>
              <Link to={arrayFind(permissions, 'menu_physical_machine') ? '/device/list' : ''}>
                <div className='home-card'>
                  <img src='assets/homepage/all_device.png' alt='' />
                  <span className='card-right'>
                    <em>{devices.data.total_devices}</em> 台<br />物理机
                </span>
                </div>
              </Link>
            </Col>
            <Col span='6'>
              <Link to={ifPermise ? '/device/setting?type=installing' : ''}>
                <div className='home-card'>
                  <img src='assets/homepage/install.png' alt='' />
                  <span className='card-right'>
                    <em>{devices.data.installing_count}</em> 台<br />正在部署
                </span>
                </div>
              </Link>
            </Col>
            <Col span='6'>
              <Link to={ifPermise ? '/device/setting?type=success' : ''}>
                <div className='home-card'>
                  <img src='assets/homepage/install_success.png' alt='' />
                  <span className='card-right'>
                    <em>{devices.data.success_count}</em> 台<br />部署成功
                </span>
                </div>
              </Link>
            </Col>
            <Col span='6'>
              <Link to={ifPermise ? '/device/setting?type=failure' : ''}>
                <div className='home-card'>
                  <img src='assets/homepage/install_fail.png' alt='' />
                  <span className='card-right'>
                    <em>{devices.data.failure_count}</em> 台<br />部署失败
                </span>
                </div>
              </Link>
            </Col>
          </Row>
          {/*<Row gutter={16} className='margin-bottom middle-card'>*/}
          {/*<Col span='8'>*/}
          {/*<Card title='机房容量统计'>*/}
          {/*<div className='pine-pos'>*/}
          {/*<Pine inspections={inspections} name='机房容量统计' id='roomChart'/>*/}
          {/*<div className='pine-center'>*/}
          {/*<p><em>-</em>台</p>*/}
          {/*<p>可分配</p>*/}
          {/*</div>*/}
          {/*</div>*/}
          {/*<div className='pine-text'>*/}
          {/*<p><span>可用率：</span>：__%</p>*/}
          {/*<p><span>总量：</span>：__</p>*/}
          {/*<p><span>可用：</span>：__</p>*/}
          {/*</div>*/}
          {/*</Card>*/}
          {/*</Col>*/}
          {/*<Col span='8'>*/}
          {/*<Card title='机架容量统计'>*/}
          {/*<div className='pine-pos'>*/}
          {/*<Pine inspections={inspections} name='机架容量统计' id='cabinetChart'/>*/}
          {/*<div className='pine-center'>*/}
          {/*<p><em>-</em>U</p>*/}
          {/*<p>可用</p>*/}
          {/*</div>*/}
          {/*</div>*/}
          {/*<div className='pine-text'>*/}
          {/*<p><span>可用率：</span>：__%</p>*/}
          {/*<p><span>总量：</span>：__</p>*/}
          {/*<p><span>可用：</span>：__</p>*/}
          {/*</div>*/}
          {/*</Card>*/}
          {/*</Col>*/}
          {/*<Col span='8'>*/}
          {/*<Card title='IP 容量统计'>*/}
          {/*<div className='pine-pos'>*/}
          {/*<Pine inspections={inspections} name='IP 容量统计' id='ipChart'/>*/}
          {/*<div className='pine-center'>*/}
          {/*<p><em>-</em>个</p>*/}
          {/*<p>可用</p>*/}
          {/*</div>*/}
          {/*</div>*/}
          {/*<div className='pine-text'>*/}
          {/*<p><span>可用率：</span>：__%</p>*/}
          {/*<p><span>总量：</span>：__</p>*/}
          {/*<p><span>可用：</span>：__</p>*/}
          {/*</div>*/}
          {/*</Card>*/}
          {/*</Col>*/}
          {/*</Row>*/}
          <Row>
            <Card
              title='近期巡检结果统计'
            >
              <RadioGroup onChange={this.onChange} defaultValue='latest_week'>
                <RadioButton value='latest_week'>近一周</RadioButton>
                <RadioButton value='latest_month'>近一月</RadioButton>
              </RadioGroup>
              <Line inspections={inspections} />
            </Card>
          </Row>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  return {
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
    data: state.get('homepage').toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}


export default connect(mapStateToProps, mapDispatchToProps)(Container);
