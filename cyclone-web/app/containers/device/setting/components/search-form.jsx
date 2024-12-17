import React from 'react';
import SearchForm from 'components/search-form';
import { DEVICE_INSTALL_STATUS, getSearchList } from "common/enums";

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'sn', label: '序列号' },
      { key: 'fixed_asset_number', label: '固资编号' },
      { key: 'category', label: '设备类型' },
      { key: 'intranet_ip', label: '内网IP' },
      { key: 'extranet_ip', label: '外网IP' },
      { key: 'server_room_name', label: '机房管理单元' },
      { key: 'server_cabinet_number', label: '机架编号' }
      // { key: 'status', label: '部署状态' }
    ],
    searchValues: {
      'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'fixed_asset_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'category': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'intranet_ip': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'extranet_ip': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'server_room_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'server_cabinet_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' }
      // 'status': { type: 'radio', list: getSearchList(DEVICE_INSTALL_STATUS) }
    }
  };
  //
  // UNSAFE_componentWillReceiveProps(props) {
  //   const { room, cabinet } = props;
  //   const searchValues = this.state.searchValues;
  //   if (!room.loading && room.loading !== this.props.room.loading) {
  //     searchValues['server_room_id'].list = room.data;
  //     this.setState({ searchValues });
  //   }
  //   if (!cabinet.loading && cabinet.loading !== this.props.cabinet.loading) {
  //     searchValues['server_cabinet_id'].list = cabinet.data;
  //     this.setState({ searchValues });
  //   }
  //
  // }


  onSearch = (values) => {
    this.props.dispatch({
      type: 'device-setting/table-data/search',
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


