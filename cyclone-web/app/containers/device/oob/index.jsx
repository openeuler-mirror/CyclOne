import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'device-oob/table-data/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'device-oob/table-data/reset'
    });
  }

  onSearch = (values) => {
    this.props.dispatch({
      type: 'device-oob/table-data/search',
      payload: {
        ...values
      }
    });
    this.props.dispatch({
      type: 'device-oob/table-data/set/selectedRows',
      payload: {
        selectedRows: [],
        selectedRowKeys: []
      }
    });
  };

  render() {
    const { tableData } = this.props.data;
    return (
      <Layout title='带外管理'>
        <SearchForm dispatch={this.props.dispatch} onSearch={this.onSearch}/>
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
    data: state.get('device-oob').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
