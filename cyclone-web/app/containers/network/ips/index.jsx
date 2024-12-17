import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';
// import { toJS } from 'utils/to-js';
import withImmutablePropsToJS from 'with-immutable-props-to-js';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'network-ips/table-data/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'network-ips/table-data/reset'
    });
  }

  render() {
    const { tableData } = this.props.data;
    return (
      <Layout title='IP分配管理'>
        <SearchForm dispatch={this.props.dispatch}/>
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
    data: state.get('network-ips'),
    userInfo: state.getIn([ 'global', 'userData' ])
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}



export default connect(mapStateToProps, mapDispatchToProps)(withImmutablePropsToJS(Container));
