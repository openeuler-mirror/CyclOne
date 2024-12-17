import React from 'react';
import SearchForm from 'components/search-form';
import { getSearchList, BUILTIN, TASK_CATEGORY, TASK_RATE, TASK_STATUS } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'title', label: '标题' },
      { key: 'builtin', label: '是否内置' },
      { key: 'category', label: '类别' },
      { key: 'rate', label: '执行频率' },
      { key: 'status', label: '状态' }
    ],
    searchValues: {
      'title': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'builtin': { type: 'radio', list: getSearchList(BUILTIN) },
      'category': { type: 'radio', list: getSearchList(TASK_CATEGORY) },
      'rate': { type: 'radio', list: getSearchList(TASK_RATE) },
      'status': { type: 'radio', list: getSearchList(TASK_STATUS) }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: 'task-list/table-data/search',
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


