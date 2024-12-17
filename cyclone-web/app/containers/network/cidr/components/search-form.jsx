import React from 'react';
import SearchForm from 'components/search-form';
import { get } from 'common/xFetch2';
import { notification } from 'antd';
import { IP_NETWORK_CATEGORY, IP_NETWORK_STATUS,getSearchList } from 'common/enums';

export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'server_room_name', label: '机房管理单元' },
      { key: 'cidr', label: '网段名称' },
      { key: 'category', label: '网段类别' },
      { key: 'switches', label: '覆盖交换机固资编号' },
      { key: 'network_area_name', label: '网络区域' }
    ],
    searchValues: {
      'server_room_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'cidr': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'category': { type: 'checkbox', list: getSearchList(IP_NETWORK_CATEGORY) },
      'switches': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'network_area_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' }
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
          if (it.id) {
            return {
              label: it.name,
              value: it.id
            };
          } else if (it.name) {
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

  //componentDidMount() {
  // this.getList('switches');
  //}


  UNSAFE_componentWillReceiveProps(props) {
    const { room, networkArea } = props;
    const searchValues = this.state.searchValues;
    // if (!room.loading && room.loading !== this.props.room.loading) {
    //   searchValues['server_room_id'].list = (room.data || []).map(it => {
    //     return {
    //       label: it.name, value: it.id
    //     };
    //   });
    //   this.setState({ searchValues });
    // }
    // if (!networkArea.loading && networkArea.loading !== this.props.networkArea.loading) {
    //   searchValues['network_area_id'].list = (networkArea.data || []).map(it => {
    //     return {
    //       label: it.name, value: it.id
    //     };
    //   });
    //   this.setState({ searchValues });
    // }
  }

  onSearch = (values) => {
    this.props.dispatch({
      type: 'network-cidr/table-data/search',
      payload: {
        ...values
      }
    });
  };

  render() {
    return (
      <SearchForm onSearch={this.onSearch} searchKeys={this.state.searchKeys} searchValues={this.state.searchValues} />
    );
  }
}


