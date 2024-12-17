import React from 'react';
import SearchForm from 'components/search-form';
import { NET_STATUS, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'server_room_name', label: '机房管理单元' },
      { key: 'name', label: '网络区域名称' },
      { key: 'physical_area', label: '关联物理区域' },
      { key: 'status', label: '状态' }
    ],
    searchValues: {
      'server_room_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'physical_area': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'status': { type: 'checkbox', list: getSearchList(NET_STATUS) }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'database-network/table-data/search',
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


