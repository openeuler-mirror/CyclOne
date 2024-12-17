import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'order-deviceCategory/table-data/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'order-deviceCategory/table-data/reset'
    });
  }

  render() {
    const { tableData, physicalArea, deviceCategory } = this.props.data;
    return (
      <Layout title='设备类型'>
        <SearchForm dispatch={this.props.dispatch}/>
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          idc={this.props.idc}
          room={this.props.room}
          physicalArea={physicalArea}
          deviceCategory={deviceCategory}
        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('order-deviceCategory').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
