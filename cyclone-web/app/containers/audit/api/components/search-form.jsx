import React from 'react';
import SearchForm from 'components/search-form';
import { API_STATUS, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'api', label: '接口地址' },
      { key: 'description', label: '接口描述' },
      { key: 'method', label: '请求方法' },
      { key: 'operator', label: '操作者' },
      { key: 'status', label: '执行状态' },
      { key: 'created_at', label: '操作时间' },
      { key: 'time', label: '耗时(s)' }
    ],
    searchValues: {
      'api': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'description': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'method': { type: 'radio', list: [ 'POST', 'PUT', 'DELETE' ] },
      'operator': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'status': { type: 'radio', list: getSearchList(API_STATUS) },
      'created_at': { type: 'rangePicker' },
      'time': { type: 'section' }
    }
  };

  onSearch = (values) => {
    if (values.created_at) {
      values.create_at_start = values.created_at[0];
      values.create_at_end = values.created_at[1];
    }
    if (values.time) {
      values.cost1 = values.time[0];
      values.cost2 = values.time[1];
    }
    delete values.created_at;
    delete values.time;
    this.props.dispatch({
      type: 'audit-api/table-data/search',
      payload: {
        ...values
      }
    });
  };

  render() {
    return (
      <SearchForm onSearch={this.onSearch} searchKeys={this.state.searchKeys} searchValues={this.state.searchValues}/>
    );
  }
}


