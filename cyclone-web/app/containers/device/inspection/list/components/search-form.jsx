import React from 'react';
import SearchForm from 'components/search-form';
import { getSearchList, RUNNING_STTAUS, INSPECTION_RESULT } from "common/enums";

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'sn', label: '序列号' },
      { key: 'fixed_asset_number', label: '固资编号' },
      { key: 'intranet_ip', label: '内网IP' },
      { key: 'running_status', label: '运行状态' },
      { key: 'health_status', label: '健康状况' },
      { key: 'createTime', label: '巡检时间' }
    ],
    searchValues: {
      'fixed_asset_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'intranet_ip': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'running_status': { type: 'radio', list: getSearchList(RUNNING_STTAUS) },
      'health_status': { type: 'radio', list: getSearchList(INSPECTION_RESULT) },
      'createTime': { type: 'rangePicker' }
    }
  };

  onSearch = (values) => {
    if (values.createTime) {
      values.start_time = values.createTime[0];
      values.end_time = values.createTime[1];
    }
    delete values.createTime;
    this.props.dispatch({
      type: 'device-inspection-list/table-data/search',
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


