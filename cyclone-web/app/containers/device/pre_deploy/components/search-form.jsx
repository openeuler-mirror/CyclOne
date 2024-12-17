import React from 'react';
import SearchForm from 'components/search-form';

//已废弃
export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'sn', label: '序列号' },
      { key: 'category', label: '设备类型' },
      { key: 'fixed_asset_number', label: '固资编号' },
      { key: 'model', label: '设备型号' },
      { key: 'usage', label: '用途' },
      { key: 'vendor', label: '厂商' }
    ],
    searchValues: {
      'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'category': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'fixed_asset_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'model': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'usage': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'vendor': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' }
    }
  };


  onSearch = (values) => {
    this.props.dispatch({
      type: 'device-pre_deploy/table-data/search',
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


