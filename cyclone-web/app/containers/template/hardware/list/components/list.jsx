import React from 'react';
import { get } from 'common/xFetch2';
import {
  Button,
  Pagination
} from 'antd';
import actions from '../actions';
import MyTable from './table';
import MyCard from './card';
import { hashHistory } from 'react-router';
import { getPermissonBtn } from 'common/utils';

class myList extends React.Component {

  state = {
    showCard: true
  };

  reload = () => {
    this.props.dispatch({
      type: 'template-hardware-list/table-data/reload'
    });
  };

  execAction = (name, record) => {
    if (name === 'copyTemplate') {
      hashHistory.push(`/template/hardware/create/${record.id}`);
    } else if (name === 'editTemplate') {
      hashHistory.push(`/template/hardware/edit/${record.id}`);
    }
    if (actions[name]) {
      actions[name]({
        record,
        reload: () => {
          this.reload();
        }
      });
    }
  };

  changePage = page => {
    this.props.dispatch({
      type: `template-hardware-list/table-data/change-page`,
      payload: {
        page
      }
    });
  };

  changePageSize = (page, pageSize) => {
    this.props.dispatch({
      type: `template-hardware-list/table-data/change-page-size`,
      payload: {
        page,
        pageSize
      }
    });
  };

  render() {
    const { tableData } = this.props;
    const { loading, pagination } = tableData;
    return (
      <div>
        <div className='operate_btns'>
          <Button
            onClick={() => hashHistory.push('template/hardware/create/new')}
            icon='plus'
            type='primary'
            disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_create')}
          >
            新增配置
          </Button>
          <span className='pull-right'>
            <Button
              onClick={() => this.setState({ showCard: true })}
              icon='appstore'
              style={{ marginRight: 8, color: this.state.showCard && '#0072ff' }}
            >
            </Button>
            <Button
              onClick={() => this.setState({ showCard: false })}
              icon='bars'
              style={{ color: !this.state.showCard && '#0072ff' }}
            >
            </Button>
          </span>
        </div>
        <div className='clearfix' />
        {
          this.state.showCard ?
            <MyCard
              dataSource={tableData.list}
              loading={loading}
              execAction={this.execAction}
              userInfo={this.props.userInfo}
            /> :
            <MyTable
              dataSource={tableData.list}
              loading={loading}
              execAction={this.execAction}
              userInfo={this.props.userInfo}
            />
        }
        <div>
          <div className='clearfix' />
          <Pagination
            showQuickJumper={true}
            showSizeChanger={true}
            current={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            onShowSizeChange={this.changePageSize}
            onChange={this.changePage}
            showTotal={(total) => `共 ${total} 条`}
          />
        </div>
      </div>
    );
  }
}

export default myList;
