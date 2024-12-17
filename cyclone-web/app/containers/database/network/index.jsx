import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'database-network/table-data/get'
    });
    //查询机房信息
    this.props.dispatch({
      type: 'database-network/room/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'database-network/table-data/reset'
    });
  }

  render() {
    const { tableData, room } = this.props.data;
    return (
      <Layout title='网络区域管理'>
        <SearchForm dispatch={this.props.dispatch} room={room} />
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          room={room}
        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('database-network').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
