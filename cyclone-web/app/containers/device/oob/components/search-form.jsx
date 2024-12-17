import React from 'react';
import SearchForm from 'components/search-form';
import { OOB_ACCESSIBLE, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'fixed_asset_number', label: '固资编号' },
      { key: 'sn', label: '序列号' },
      { key: 'ip', label: 'IP' },
      { key: 'intranet_ip', label: '内网IP' },
      { key: 'oob_accessible', label: '带外状态' }
    ],
    searchValues: {
      'fixed_asset_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'ip': { type: 'input', placeholder: '模糊查询所有IP,支持英文逗号,空格,换行分隔' },
      'intranet_ip': { type: 'input', placeholder: '精确查询内网IP,支持英文逗号,空格,换行分隔' },
      'oob_accessible': { type: 'checkbox', list: getSearchList(OOB_ACCESSIBLE) }
    }
  };

  render() {
    return (
      <SearchForm onSearch={this.props.onSearch} searchKeys={this.state.searchKeys} searchValues={this.state.searchValues} />
    );
  }
}


