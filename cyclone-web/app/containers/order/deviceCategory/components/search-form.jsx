import React from 'react';
import SearchForm from 'components/search-form';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'category', label: '设备类型' },
      { key: 'hardware', label: '硬件配置' },
      { key: 'remark', label: '备注' }
    ],
    searchValues: {
      'category': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'hardware': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'remark': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'order-deviceCategory/table-data/search',
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


