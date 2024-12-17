import React from 'react';
import SearchForm from 'components/search-form';
import { YES_NO, NETWORK_IPS_SCOPE, IP_NETWORK_CATEGORY, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'ip', label: 'IP地址' },
      { key: 'cidr', label: '网段名称' },
      { key: 'sn', label: '序列号' },
      { key: 'fixed_asset_number', label: '固资编号' },
      { key: 'is_used', label: '是否被使用' },
      { key: 'scope', label: 'IP作用范围' },
      { key: 'ipnetwork_category', label: '类别' }
    ],
    searchValues: {
      'ip': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'cidr': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'fixed_asset_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'is_used': { type: 'checkbox', list: getSearchList(YES_NO) },
      'scope': { type: 'checkbox', list: getSearchList(NETWORK_IPS_SCOPE) },
      'ipnetwork_category': { type: 'checkbox', list: getSearchList(IP_NETWORK_CATEGORY) }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'network-ips/table-data/search',
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


