import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'database-usite/table-data/get'
    });
    //获取机房数据
    this.props.dispatch({
      type: 'database-usite/room/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'database-usite/table-data/reset'
    });
  }

  render() {
    const { tableData, room } = this.props.data;
    return (
      <Layout title='机位信息管理'>
        <SearchForm
          dispatch={this.props.dispatch}
          room={room}
          idc={this.props.idc}
        />
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
    data: state.get('database-usite').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
    idc: state.getIn([ 'global', 'idc' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
