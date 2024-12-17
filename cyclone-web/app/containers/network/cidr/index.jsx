import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'network-cidr/table-data/get'
    });
    //查询机房信息
    this.props.dispatch({
      type: 'network-cidr/room/get'
    });
    //查询网络区域
    this.props.dispatch({
      type: 'network-cidr/networkArea/get'
    });
    //查询网络设备信息
    this.props.dispatch({
      type: 'network-cidr/device/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'network-cidr/table-data/reset'
    });
  }

  render() {
    const { tableData, room, device, networkArea } = this.props.data;
    return (
      <Layout title='IP网段管理'>
        <SearchForm dispatch={this.props.dispatch} room={room} networkArea={networkArea} />
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          room={room}
          device={device}
        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('network-cidr').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
