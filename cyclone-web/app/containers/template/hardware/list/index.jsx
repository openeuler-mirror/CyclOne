import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import List from './components/list';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'template-hardware-list/table-data/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'template-hardware-list/table-data/reset'
    });
  }

  render() {
    const { tableData } = this.props.data;
    return (
      <Layout title='硬件配置'>
        <SearchForm
          dispatch={this.props.dispatch}
        />
        <List
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
    data: state.get('template-hardware-list').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
