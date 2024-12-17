import React from 'react';
import SearchForm from 'components/search-form';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'http_method', label: '请求方式' },
      { key: 'url', label: '路由' },
      { key: 'category_name', label: '操作类型' }
    ],
    searchValues: {
      'http_method': { type: 'radio', list: [ 'GET', 'PUT', 'POST', 'DELETE' ] },
      'url': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'category_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'audit-log/table-data/search',
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


