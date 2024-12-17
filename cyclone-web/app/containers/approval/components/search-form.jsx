import React from 'react';
import SearchForm from 'components/search-form';
import { APPROVAL_STATUS, APPROVAL_TYPE, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    pedingSearchKeys: [{ key: 'type', label: '审批类型' }],
    searchKeys: [
      { key: 'type', label: '审批类型' },
      { key: 'status', label: '审批状态' }
    ],
    searchValues: {
      'type': { type: 'radio', list: getSearchList(APPROVAL_TYPE) },
      'status': { type: 'radio', list: getSearchList(APPROVAL_STATUS) }
    }
  };

  onSearch = (values) => {
    this.props.dispatch({
      type: `approval/${this.props.type}-table-data/search`,
      payload: {
        ...values
      }
    });
  };

  // getInitialValue = () => {
  //   const data = this.props.initialValue;
  //   const list = [];
  //   if (data && JSON.stringify(data) == '{}') {
  //     return; 
  //   }
  //   for (let item in data) {
  //     list.push({
  //       key: item, keyLabel: '', valueLabel: '', value: data[item]
  //     });
  //   }
  //   return list;
  // }

  render() {
    let searchKeys = [];
    if (this.props.type == 'pending') {
      searchKeys = this.state.pedingSearchKeys;
    } else {
      searchKeys = this.state.searchKeys;
    }
    return (
      <SearchForm onSearch={this.onSearch} searchKeys={searchKeys} searchValues={this.state.searchValues}/>
    );
  }
}


