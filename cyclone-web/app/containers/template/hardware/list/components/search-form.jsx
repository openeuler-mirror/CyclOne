import React from 'react';
import SearchForm from 'components/search-form';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'name', label: '名称' },
      // { key: 'vendor', label: '厂商' },
      // { key: 'model_name', label: '型号' },
      { key: 'builtin', label: '配置类型' }
    ],
    searchValues: {
      'name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      // 'vendor': { type: 'input' },
      // 'model_name': { type: 'input' },
      'builtin': { type: 'radio', list: [
        { value: 'Yes', label: '内置' },
        { value: 'No', label: '自定义' }
      ] }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'template-hardware-list/table-data/search',
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


