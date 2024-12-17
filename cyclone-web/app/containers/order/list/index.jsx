import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'order-list/table-data/get'
    });
    this.props.dispatch({
      type: 'order-list/physical-area/get'
    });
    this.props.dispatch({
      type: 'order-list/device-categories/get'
    });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'order-list/table-data/reset'
    });
  }

  render() {
    const { tableData, physicalArea, deviceCategory } = this.props.data;
    return (
      <Layout title='订单管理'>
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
    data: state.get('order-list').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
    dict: state.getIn([ 'global', 'dict' ]).toJS(),
    room: state.getIn([ 'global', 'room' ]).toJS(),
    idc: state.getIn([ 'global', 'idc' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
