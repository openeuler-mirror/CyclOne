import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'database-store/table-data/get'
    });
    //获取数据中心数据
    this.props.dispatch({
      type: 'database-store/idc/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'database-store/table-data/reset'
    });
  }

  render() {
    const { tableData, idc } = this.props.data;
    return (
      <Layout title='库房信息管理'>
        <SearchForm dispatch={this.props.dispatch}/>
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
    data: state.get('database-store').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
