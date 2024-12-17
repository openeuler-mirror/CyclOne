import React from 'react';
import SearchForm from 'components/search-form';
import { CAB_TYPE, CAB_STATUS, YES_NO, getSearchList } from 'common/enums';
import { get } from 'common/xFetch2';
import { notification } from 'antd';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'number', label: '机架编号' },
      { key: 'server_room_name', label: '机房管理单元' },
      { key: 'idc_id', label: '数据中心' },
      { key: 'network_area_name', label: '网络区域' },
      { key: 'type', label: '类型' },
      { key: 'status', label: '机架状态' },
      { key: 'is_enabled', label: '是否启用' },
      { key: 'is_powered', label: '是否开电' }
    ],
    searchValues: {
      'number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'server_room_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'idc_id': { type: 'checkbox', list: [] },
      'network_area_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'type': { type: 'checkbox', list: getSearchList(CAB_TYPE) },
      'status': { type: 'checkbox', list: getSearchList(CAB_STATUS) },
      'is_enabled': { type: 'checkbox', list: getSearchList(YES_NO) },
      'is_powered': { type: 'checkbox', list: getSearchList(YES_NO) }
    }
  };

  getListId = (name) => {
    const searchValues = this.state.searchValues;
    get(`/api/cloudboot/v1/idcs`, { page: 1, page_size: 100 }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      const list = res.content.records || [];
      if (list.length > 0) {
        searchValues[name].list = list.map(it => {
          return {
            label: it.name,
            value: it.id
          };
        });
      }
    });
  };

  componentDidMount() {
    this.getListId('idc_id');
  }

  UNSAFE_componentWillReceiveProps(props) {
    const { network, room, idc } = props;
    const searchValues = this.state.searchValues;
    // if (!network.loading && network.loading !== this.props.network.loading) {
    //   searchValues['network_area_id'].list = network.data.map(it => {
    //     return {
    //       label: it.name,
    //       value: it.id
    //     };
    //   });
    //   this.setState({ searchValues });
    // }
    // if (!room.loading && room.loading !== this.props.room.loading) {
    //   searchValues['server_room_id'].list = room.data;
    //   this.setState({ searchValues });
    // }
    // if (!idc.loading && idc.loading !== this.props.idc.loading) {
    //   searchValues['idc_id'].list = idc.data;
    //   this.setState({ searchValues });
    // }
  }


  onSearch = (values) => {
    this.props.dispatch({
      type: 'database-cabinet/table-data/search',
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


