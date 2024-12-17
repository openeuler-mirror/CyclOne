import React from 'react';
import SearchForm from 'components/search-form';
import { OPERATION_STATUS, OOB_ACCESSIBLE, getSearchList } from 'common/enums';
import { get } from 'common/xFetch2';
import { notification } from 'antd';

export default class Search extends React.Component {


  state = {
    searchKeys: [
      { key: 'sn', label: '序列号' },
      { key: 'category', label: '设备类型列表' },//枚举
      { key: 'category_pre_deploy', label: '待部署设备类型' },//枚举
      { key: 'fixed_asset_number', label: '固资编号' },
      { key: 'model', label: '设备型号' }, //枚举
      { key: 'operation_status', label: '运营状态' },
      { key: 'usage', label: '用途' }, //枚举
      { key: 'vendor', label: '厂商' }, //枚举
      { key: 'ip', label: 'IP' },
      { key: 'intranet_ip', label: '内网IP' },
      { key: 'extranet_ip', label: '外网IP' },
      { key: 'idc_id', label: '数据中心' }, //枚举
      { key: 'server_room_name', label: '机房管理单元' }, //枚举
      { key: 'physical_area', label: '物理区域' }, //枚举
      { key: 'server_cabinet_number', label: '机架编号' },
      { key: 'server_usite_number', label: '机位编号' },
      { key: 'oob_accessible', label: '带外状态' }
    ],
    searchValues: {
      'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'category': { type: 'checkbox', list: [] },
      'category_pre_deploy': { type: 'checkbox', list: [] },
      'fixed_asset_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'model': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'operation_status': { type: 'radio', list: getSearchList(OPERATION_STATUS) },
      'usage': { type: 'checkbox', list: [] },
      'vendor': { type: 'checkbox', list: [] },
      'ip': { type: 'input', placeholder: '模糊匹配内外网以及带外IP,支持英文逗号,空格,换行分隔' },
      'intranet_ip': { type: 'input', placeholder: '精确匹配内网IP,支持英文逗号,空格,换行分隔' },
      'extranet_ip': { type: 'input', placeholder: '精确匹配外网IP,支持英文逗号,空格,换行分隔' },
      'idc_id': { type: 'checkbox', list: [] },
      'server_room_name': { type: 'input' },
      'physical_area': { type: 'checkbox', list: [] },
      'server_cabinet_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'server_usite_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'oob_accessible': { type: 'checkbox', list: getSearchList(OOB_ACCESSIBLE) }

    }
  };

  getList = (name) => {
    const searchValues = this.state.searchValues;
    get(`/api/cloudboot/v1/devices/query-params/${name}`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      const list = res.content.list || [];
      if (list.length > 0 && searchValues[name]) {
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

  getListId = (urlParams, name) => {
    const searchValues = this.state.searchValues;
    get(`/api/cloudboot/v1/devices/query-params/${urlParams}`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      const list = res.content.list || [];
      if (list.length > 0 && searchValues[name]) {
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
    this.getList('vendor');
    this.getList('physical_area');
    this.getList('usage');
    this.getList('category');
    this.getList('category_pre_deploy');
    this.getListId('idc', 'idc_id');
    this.getListId('server_room', 'server_room_id');
  }

  // UNSAFE_componentWillReceiveProps(props) {
  //   const { idc, room } = props;
  //   const searchValues = this.state.searchValues;
  //   if (!idc.loading && idc.loading !== this.props.idc.loading) {
  //     searchValues['idc_id'].list = idc.data;
  //   }
  //   if (!room.loading && room.loading !== this.props.room.loading) {
  //     searchValues['server_room_id'].list = room.data;
  //   }
  //   this.setState({ searchValues });
  // }

  render() {
    return (
      <SearchForm onSearch={this.props.onSearch} searchKeys={this.state.searchKeys} searchValues={this.state.searchValues} />
    );
  }
}


