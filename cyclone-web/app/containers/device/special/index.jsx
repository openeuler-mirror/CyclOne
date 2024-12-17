import React from 'react';
import { connect } from 'react-redux';
import SearchForm from './components/search-form';
import Table from './components/table';
import Layout from 'components/layout/page-layout';

class Container extends React.Component {
  componentDidMount() {
    this.props.dispatch({
      type: 'device-special/table-data/get'
    });
    // this.props.dispatch({
    //   type: 'device-special/idc/get'
    // });
    // this.props.dispatch({
    //   type: 'device-special/room/get'
    // });
    // this.props.dispatch({
    //   type: 'device-special/networkArea/get'
    // });
  }

  componentWillUnmount() {
    this.props.dispatch({
      type: 'device-special/table-data/reset'
    });
  }

  // onSearch = (values) => {
  //   this.props.dispatch({
  //     type: 'device-special/table-data/search',
  //     payload: {
  //       ...values
  //     }
  //   });
  // };

  render() {
    const { tableData } = this.props.data;
    return (
      <Layout title='特殊设备'>
        <SearchForm
          dispatch={this.props.dispatch}
        />
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          liableUser={this.props.dict.liableUser}
          room={this.props.room}

        />
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('device-special').toJS(),
    userInfo: state.getIn(['global', 'userData']).toJS(),
    dict: state.getIn(['global', 'dict']).toJS(),
    room: state.getIn(['global', 'room']).toJS(),
    networkArea: state.getIn(['global', 'networkArea']).toJS(),
    idc: state.getIn(['global', 'idc']).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
