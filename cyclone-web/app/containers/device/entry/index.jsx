import React from 'react';
import { connect } from 'react-redux';
import Table from './components/table';
import Layout from 'components/layout/page-layout';
import { getBreadcrumb } from 'common/utils';

class Container extends React.Component {

  //type: 'recycle' 回收重装
  //form: "approval" 审批管理
  constructor(props) {
    super(props);
    this.state = {
      from: 'origin',
      ...props.location.query
    };
  }

  componentDidMount() {
    console.log(this.state);
  }



  render() {
    const dataSource = this.props.location.state || [];
    if (dataSource.length > 0) {
      dataSource.map(item => {
        item.server_room_name = item.server_room.name;
        item.need_extranet_ip = 'no'; //给一个默认值
        item.need_intranet_ipv6 = 'no'; //给一个默认值
        item.need_extranet_ipv6 = 'no'; //给一个默认值
      });
    }
    let name = '上架部署';
    if (this.state.from === 'approval') {
      if (this.state.type === 'recycle') {
        name = '回收重装';
      } else {
        name = '物理机重装';
      }
    }
    return (
      <Layout>
        <div style={{ marginTop: -10 }}>
          {getBreadcrumb(name)}
        </div>
        <Table from={this.state.from} type={this.state.type} dataSource={dataSource} ip={this.props.ip} inipv6={this.props.inipv6} exipv6={this.props.exipv6} hardwareData={this.props.hardwareData} sysData={this.props.sysData} dispatch={this.props.dispatch} userInfo={this.props.userInfo} />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
    hardwareData: state.getIn([ 'device-entry', 'hardwareData' ]).toJS(),
    ip: state.getIn([ 'device-entry', 'ip' ]),
    inipv6: state.getIn([ 'device-entry', 'inipv6' ]),
    exipv6: state.getIn([ 'device-entry', 'exipv6' ]),
    sysData: state.getIn([ 'device-entry', 'sysData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
