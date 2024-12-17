import React from 'react';
import SearchForm from 'components/search-form';
import { IDC_STATUS, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'name', label: '库房管理单元' },
      { key: 'store_room_manager', label: '库房负责人' },
      { key: 'vendor_manager', label: '供应商负责人' }
      // { key: 'status', label: '状态' }
    ],
    searchValues: {
      'name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'store_room_manager': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'vendor_manager': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' }
      // 'status': { type: 'checkbox', list: getSearchList(IDC_STATUS) }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'database-store/table-data/search',
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


