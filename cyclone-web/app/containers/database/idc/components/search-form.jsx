import React from 'react';
import SearchForm from 'components/search-form';
import { IDC_USAGE, IDC_STATUS, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'name', label: '数据中心' },
      { key: 'usage', label: '用途' },
      { key: 'first_server_room', label: '一级机房' },
      { key: 'vendor', label: '供应商' },
      { key: 'status', label: '状态' }
    ],
    searchValues: {
      'name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'usage': { type: 'checkbox', list: getSearchList(IDC_USAGE) },
      'first_server_room': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'vendor': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'status': { type: 'checkbox', list: getSearchList(IDC_STATUS) }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'database-idc/table-data/search',
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


