import React from 'react';
import SearchForm from 'components/search-form';
import { DEVICE_SETTING_RULE_CATEGORY, getSearchList } from "common/enums";

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'rule_category', label: '规则类别' },
    ],
    searchValues: {
      'rule_category': { type: 'checkbox', list: getSearchList(DEVICE_SETTING_RULE_CATEGORY) },
    }
  };
 
  onSearch = (values) => {
    this.props.dispatch({
      type: 'device-setting-rules/table-data/search',
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


