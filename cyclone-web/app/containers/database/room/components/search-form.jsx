import React from 'react';
import SearchForm from 'components/search-form';
import { IDC_STATUS, getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'idc_id/first_server_room', label: '数据中心 / 所属一级机房' },
      { key: 'name', label: '机房管理单元' },
      { key: 'city', label: '城市' },
      { key: 'address', label: '机房地址' },
      { key: 'server_room_manager', label: '机房负责人' },
      { key: 'vendor_manager', label: '供应商负责人' },
      { key: 'network_asset_manager', label: '网络资产负责人' },
      { key: 'support_phone_number', label: '7*24小时保障电话' },
      { key: 'status', label: '状态' }
    ],
    searchValues: {
      'idc_id/first_server_room': { type: 'cascader_ysnc', list: [] },
      'name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'city': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'address': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'server_room_manager': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'vendor_manager': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'network_asset_manager': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'support_phone_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'status': { type: 'checkbox', list: getSearchList(IDC_STATUS) }
    }
  };

  UNSAFE_componentWillReceiveProps(props) {
    const { idc } = props;
    const searchValues = this.state.searchValues;
    if (!idc.loading && idc.loading !== this.props.idc.loading) {
      searchValues['idc_id/first_server_room'].list = (idc.data || []).map(it => {
        return {
          label: it.name, value: it.id, children: it.first_server_room
        };
      });
      this.setState({ searchValues });
    }
  }

  onSearch = (values) => {
    if (values['idc_id/first_server_room']) {
      values['idc_id'] = values['idc_id/first_server_room'][0];
      values['first_server_room'] = values['idc_id/first_server_room'][1];
    }
    delete values['idc_id/first_server_room'];
    this.props.dispatch({
      type: 'database-room/table-data/search',
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


