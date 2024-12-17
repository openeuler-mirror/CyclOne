import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'network-device/table-data/get'
    });
    this.props.dispatch({
      type: 'network-device/idc/get'
    });
    this.props.dispatch({
      type: 'network-device/room/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'network-device/table-data/reset'
    });
  }

  render() {
    const { tableData, idc, room } = this.props.data;
    return (
      <Layout title='网络设备管理'>
        <SearchForm dispatch={this.props.dispatch} idc={idc} room={room}/>
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          idc={idc}
        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('network-device').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
