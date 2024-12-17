import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'device-list/table-data/get'
    });
    // this.props.dispatch({
    //   type: 'device-list/idc/get'
    // });
    // this.props.dispatch({
    //   type: 'device-list/room/get'
    // });
    // this.props.dispatch({
    //   type: 'device-list/networkArea/get'
    // });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'device-list/table-data/reset'
    });
  }

  onSearch = (values) => {
    this.props.dispatch({
      type: 'device-list/table-data/search',
      payload: {
        ...values
      }
    });
  };

  render() {
    const { tableData } = this.props.data;
    return (
      <Layout title='物理机列表'>
        <SearchForm networkArea={this.props.networkArea} room={this.props.room} idc={this.props.idc} dispatch={this.props.dispatch} onSearch={this.onSearch}/>
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          liableUser={this.props.dict.liableUser}
        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('device-list').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS(),
    dict: state.getIn([ 'global', 'dict' ]).toJS(),
    room: state.getIn([ 'global', 'room' ]).toJS(),
    networkArea: state.getIn([ 'global', 'networkArea' ]).toJS(),
    idc: state.getIn([ 'global', 'idc' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
