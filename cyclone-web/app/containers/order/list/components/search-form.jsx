import React from 'react';
import SearchForm from 'components/search-form';
import { ORDER_STATUS, getSearchList } from 'common/enums';
import { get } from 'common/xFetch2';
import { notification } from 'antd';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'number', label: '订单编号' },
      { key: 'status', label: '状态' },
      { key: 'usage', label: '用途' },
      { key: 'physical_area', label: '物理区域' }
    ],
    searchValues: {
      'number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'status': { type: 'checkbox', list: getSearchList(ORDER_STATUS) },
      'usage': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'physical_area': { type: 'checkbox', list: [] }
    }
  };

  getList = (name) => {
    const searchValues = this.state.searchValues;
    get(`/api/cloudboot/v1/devices/query-params/${name}`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      const list = res.content.list || [];
      if (list.length > 0) {
        searchValues[name].list = list.map(it => {
          if (it.name) {
            return {
              label: it.name,
              value: it.name
            };
          } else {
            return {
              label: 'blank',
              value: 'blank'
            };
          }
        });
      }
    });
  };

  componentDidMount() {
    this.getList('physical_area');
  }

  onSearch = (values) => {
    this.props.dispatch({
      type: 'order-list/table-data/search',
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


