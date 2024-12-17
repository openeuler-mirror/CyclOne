import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'database-cabinet/table-data/get'
    });
    // 获取网络区域信息
    this.props.dispatch({
      type: 'database-cabinet/network/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'database-cabinet/table-data/reset'
    });
  }

  render() {
    const { tableData, network } = this.props.data;
    return (
      <Layout title='机架信息管理'>
        <SearchForm room={this.props.room}
          dispatch={this.props.dispatch}
          network={network}
          idc={this.props.idc}

        />
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          network={network}
        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('database-cabinet').toJS(),
    room: state.getIn([ 'global', 'room' ]).toJS(),
    idc: state.getIn([ 'global', 'idc' ]).toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
