import React from 'react';
import SearchForm from 'components/search-form';


export default class Search extends React.Component {
  state = {
    searchKeys: [
      { key: 'fixed_asset_number', label: '固资编号' },
      { key: 'sn', label: '序列号' },
      { key: 'name', label: '名称' },
      { key: 'idc_id', label: '数据中心' },
      { key: 'server_cabinet_number', label: '机架编号' },
      { key: 'server_room_name', label: '机房管理单元' },
      { key: 'tor', label: 'TOR' },
      { key: 'model', label: '产品型号' },
      { key: 'vendor', label: '厂商' },
      { key: 'os', label: '操作系统' },
      { key: 'usage', label: '用途' },
      { key: 'status', label: '状态' }
    ],
    searchValues: {
      'fixed_asset_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'idc_id': { type: 'checkbox', list: [] },
      'server_cabinet_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'server_room_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'tor': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'model': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'vendor': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'os': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'usage': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
      'status': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' }
    }
  };

  UNSAFE_componentWillReceiveProps(props) {
    const { idc, room } = props;
    const searchValues = this.state.searchValues;
    if (!idc.loading && idc.loading !== this.props.idc.loading) {
      searchValues['idc_id'].list = (idc.data || []).map(it => {
        return {
          label: it.name, value: it.id
        };
      });
      this.setState({ searchValues });
    }
  }

  onSearch = (values) => {
    if (values.tor && values.tor.indexOf('+') != -1) {
      values.tor = encodeURIComponent(values.tor); 
    }
    this.props.dispatch({
      type: 'network-device/table-data/search',
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


