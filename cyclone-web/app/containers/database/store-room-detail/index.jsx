import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout';
import { getBreadcrumb } from 'common/utils';
import Detail from './components/form-detail';
import Table from './components/table';
import { getPermissonBtn } from 'common/utils';

class Container extends React.Component {
  state = {
    id: this.props.params.id
  };

  componentDidMount() {
    this.reload();
    this.props.dispatch({
      type: 'database-store-detail/detail-info/get',
      payload: this.state.id
    });
  }

  //获取虚拟货架列表
  reload = () => {
    this.props.dispatch({
      type: 'database-store-detail/table-data/search',
      payload: { store_room_id: this.state.id }
    });
  };

  render() {
    const { detailInfo, tableData } = this.props.data;

    return (
      <Layout>
        {getBreadcrumb(this.state.id, '/database/store-room')}
        <h3 className='detail-title'>基本信息</h3>
        <div className='detail-info'>
          <Detail detailInfo={detailInfo}/>
        </div>
        <h3 className='detail-title'>虚拟货架</h3>
        <Table
          dispatch={this.props.dispatch}
          tableData={tableData}
          userInfo={this.props.userInfo}
          reload={this.reload}
          id={this.state.id}
        />
      </Layout>
    );
  }
}
function mapStateToProps(state) {
  return {
    data: state.get('database-store-detail').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
