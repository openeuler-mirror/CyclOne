import React from 'react';
import { connect } from 'react-redux';
import MirrorTable from './components/mirror-table';
import SysTable from './components/sys-table';
import Layout from 'components/layout/page-layout-tabs';
import { Tabs } from 'antd';
const TabPane = Tabs.TabPane;
import { hashHistory } from 'react-router';


class Container extends React.Component {
  componentDidMount() {
    //获取操作系统族数据
    this.props.dispatch({
      type: 'template-system/osFamily/get'
    });
  }

  renderSystemContent = () => {
    const { systemConfig, osFamily } = this.props.data;
    return (
      <div>
        <SysTable
          dispatch={this.props.dispatch}
          tableData={systemConfig}
          userInfo={this.props.userInfo}
          osFamily={osFamily}
        />
      </div>
    );
  };
  renderMirrorContent = () => {
    const { mirrorInstallTpl, osFamily } = this.props.data;
    return (
      <div>
        <MirrorTable
          dispatch={this.props.dispatch}
          tableData={mirrorInstallTpl}
          userInfo={this.props.userInfo}
          osFamily={osFamily}
        />
      </div>
    );
  };

  onTabChange = (key) => {
    const currentQuery = this.props.location.query;
    currentQuery.type = key;
    hashHistory.push(this.props.location);
  };

  render() {
    const { query } = this.props.location;
    const defaultKey = query.type ? query.type : 'systemConfig';
    return (
      <Layout title='装机配置'>
        <Tabs onChange={this.onTabChange} activeKey={defaultKey} defaultActiveKey={defaultKey} type='card'>
          <TabPane tab='PXE配置' key='systemConfig'>
            {this.renderSystemContent()}
          </TabPane>
          <TabPane tab='镜像配置' key='mirrorInstallTpl'>
            {this.renderMirrorContent()}
          </TabPane>
        </Tabs>
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('template-system').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
