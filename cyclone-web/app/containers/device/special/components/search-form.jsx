import React from 'react';
import SearchForm from 'components/search-form';
import { OPERATION_STATUS, getSearchList } from 'common/enums';
import { get } from 'common/xFetch2';
import { notification } from 'antd';

export default class Search extends React.Component {

  state = {
      searchKeys: [
            { key: 'sn', label: '序列号' },
            { key: 'ip', label: 'IP' },
            { key: 'model', label: '设备型号' }, //枚举
            { key: 'vendor', label: '厂商' }, //枚举
            { key: 'server_room_name', label: '机房管理单元' },
            { key: 'server_cabinet_number', label: '机架编号' },
            { key: 'server_usite_number', label: '机位编号' },
            { key: 'hardware_remark', label: '硬件说明' }
        ],
      searchValues: {
          'sn': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
          'ip': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
          'model': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
          'vendor': { type: 'checkbox', list: [] },
          'server_room_name': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
          'server_usite_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
          'server_cabinet_number': { type: 'input', placeholder: '查询条件支持英文逗号,空格,换行分隔' },
          'hardware_remark': { type: 'input', placeholder: '查询条件支持模糊查询' }
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
      this.getList('vendor');
    }

  onSearch = (values) => {
      this.props.dispatch({
          type: `device-special/table-data/search`,
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


