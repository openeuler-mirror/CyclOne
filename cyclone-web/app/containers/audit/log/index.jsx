import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'audit-log/table-data/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'audit-log/table-data/reset'
    });
  }

  render() {
    const { tableData } = this.props.data;
    return (
      <Layout title='操作记录'>
        <SearchForm dispatch={this.props.dispatch} />
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('audit-log').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
