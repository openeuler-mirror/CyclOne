import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout-tabs';
import { Tabs } from 'antd';
const TabPane = Tabs.TabPane;
import { hashHistory } from 'react-router';


class Container extends React.Component {

  componentDidMount() {
    const type = this.props.location.query.type;
    //存type
    this.props.dispatch({
      type: 'device-setting-rules/type/get',
      payload: type === 'all' ? null : type
    });
    this.props.dispatch({
      type: 'device-setting-rules/table-data/get',
      payload: type === 'all' ? null : type
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'device-setting-rules/table-data/reset'
    });
  }

  renderContent = (rule_category) => {
    const { tableData } = this.props.data;
    return (
      <div>
        <SearchForm
          dispatch={this.props.dispatch}
        />
        <Table
          rule_category={rule_category}
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
        />
      </div>
    );
  };
  getContent = (key) => {
    //存type
    this.props.dispatch({
      type: 'device-setting-rules/type/get',
      payload: key === 'all' ? null : key
    });
    const currentQuery = this.props.location.query;
    currentQuery.type = key;
    hashHistory.push(this.props.location);
    this.props.dispatch({
      type: 'device-setting-rules/table-data/search',
      payload: {
        rule_category: key === 'all' ? null : key
      }
    });
  };

  render() {
    const { query } = this.props.location;
    const defaultKey = query.type ? query.type : 'all';
    return (
      <Layout title='规则列表'>
        <Tabs onChange={this.getContent} activeKey={defaultKey} defaultActiveKey={defaultKey} type='card'>
          <TabPane tab='全部' key='all'>{this.renderContent('all')}</TabPane>
          <TabPane tab='操作系统' key='os'>{this.renderContent('os')}</TabPane>
          <TabPane tab='阵列结构' key='raid'>{this.renderContent('raid')}</TabPane>
          <TabPane tab='网络配置' key='network'>{this.renderContent('network')}</TabPane>
        </Tabs>
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('device-setting-rules').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
