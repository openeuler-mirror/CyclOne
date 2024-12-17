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
    //存type，否则关闭查看日志模态框刷新时找不到type
    this.props.dispatch({
      type: 'device-setting/type/get',
      payload: type === 'all' ? null : type
    });
    this.props.dispatch({
      type: 'device-setting/table-data/get',
      payload: type === 'all' ? null : type
    });
    this.loadNum();
  }


  componentWillUnmount() {
    this.props.dispatch({
      type: 'device-setting/table-data/reset'
    });
  }

  loadNum = () => {
    //获取数目
    this.props.dispatch({
      type: 'device-setting/statistics/get'
    });
  };
  renderContent = (status) => {
    //console.log(this.props)
    const { tableData } = this.props.data;
    return (
      <div>
        <SearchForm
          dispatch={this.props.dispatch}
          room={this.props.room}
          cabinet={this.props.cabinet}
        />
        <Table
          status={status}
          loadNum={this.loadNum}
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
      type: 'device-setting/type/get',
      payload: key === 'all' ? null : key
    });
    const currentQuery = this.props.location.query;
    currentQuery.type = key;
    hashHistory.push(this.props.location);
    this.props.dispatch({
      type: 'device-setting/table-data/search',
      payload: {
        status: key === 'all' ? null : key
      }
    });
    this.loadNum();
  };

  renderTab = (num, title) => {
    if (num > 0) {
      return <span>{title}<span className='spanCircle fill'>{num}</span></span>;
    }
    return title;
  };

  render() {
    const { statistics } = this.props.data;
    const statusNum = statistics.data || {};
    const { query } = this.props.location;
    const defaultKey = query.type ? query.type : 'all';
    return (
      <Layout title='部署列表'>
        <Tabs onChange={this.getContent} activeKey={defaultKey} defaultActiveKey={defaultKey} type='card'>
          <TabPane tab='全部' key='all'>{this.renderContent('all')}</TabPane>
          <TabPane tab={this.renderTab(statusNum.preinstall_count, '等待部署')} key='pre_install'>{this.renderContent('pre_install')}</TabPane>
          <TabPane tab={this.renderTab(statusNum.installing_count, '正在部署')} key='installing'>{this.renderContent('installing')}</TabPane>
          <TabPane tab={this.renderTab(statusNum.success_count, '部署成功')} key='success'>{this.renderContent('success')}</TabPane>
          <TabPane tab={this.renderTab(statusNum.failure_count, '部署失败')} key='failure'>{this.renderContent('failure')}</TabPane>
        </Tabs>
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('device-setting').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
    room: state.getIn([ 'global', 'room' ]).toJS(),
    cabinet: state.getIn([ 'global', 'cabinet' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
