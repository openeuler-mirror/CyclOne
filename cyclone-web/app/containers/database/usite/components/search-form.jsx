import React from 'react';
import SearchForm from 'components/search-form';
import { USITE_STATUS, USITE_PORT_RATE, getSearchList } from 'common/enums';
import { get } from 'common/xFetch2';
import { notification } from 'antd';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'usite_number', label: '机位编号' },
      { key: 'idc_id', label: '数据中心' },
      { key: 'server_room_name', label: '机房管理单元' },
      { key: 'cabinet_number', label: '机架编号' },
      { key: 'height', label: '机位高度' },
      { key: 'status', label: '机位状态' },
      { key: 'physical_area', label: '物理区域' },
      { key: 'la_wa_port_rate', label: '内外网端口速率' },
    ],
    searchValues: {
      'usite_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'idc_id': { type: 'checkbox', list: [] },
      'server_room_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'cabinet_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'height': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'status': { type: 'checkbox', list: getSearchList(USITE_STATUS) },
      'physical_area': { type: 'checkbox', list: [] },
      'la_wa_port_rate': { type: 'checkbox', list: getSearchList(USITE_PORT_RATE) },
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
    this.getList('physical_area');
    this.getListId('idc_id');
  }

  UNSAFE_componentWillReceiveProps(props) {
    const { room, idc } = props;
    const searchValues = this.state.searchValues;
    // if (!room.loading && room.loading !== this.props.room.loading) {
    //   searchValues['server_room_id'].list = (room.data || []).map(it => {
    //     return {
    //       label: it.name, value: it.id
    //     };
    //   });
    //   this.setState({ searchValues });
    // }
    // if (!idc.loading && idc.loading !== this.props.idc.loading) {
    //   searchValues['idc_id'].list = idc.data;
    //   this.setState({ searchValues });
    // }
  }

  onSearch = (values) => {
    this.props.dispatch({
      type: 'database-usite/table-data/search',
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


