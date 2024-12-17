import React from 'react';
import {
  Row
} from 'antd';
import handleAction from './sync-actions/index';
import asyncActions from './async-actions/index';
import M from 'immutable';
import { renderFormDetail } from 'common/utils';
import { OPERATION_STATUS, BUILTIN, DEVICE_MAINTENANCE_SERVICE_STATUS } from 'common/enums';


export default class OperationTarget extends React.Component {
  componentDidMount() {
    if (!this.props.data) {
      this.dispatch({
        type: 'device/detail/get',
        payload: this.props.sn
      });
    }
  }

  constructor(props) {
    super(props);
    this.state = {
      data: M.fromJS({
        detailInfo: {
          loading: true,
          data: {
            device_page_resp: {

            },
            //device_lifecycle_detail_page: {
//
            //}
          }
        }
      })
    };
  }

  getState = () => {
    return this.state.data;
  };
  dispatch = action => {
    if (asyncActions[action.type]) {
      asyncActions[action.type](
        this.state.data,
        action,
        this.dispatch,
        this.getState
      );
      return;
    }

    const data = handleAction(this.state.data, action);
    this.setState({
      data
    });
  };
  render() {
    let data = {};
    //let lifecycleData = {};
    //console.log("show props", this.props);
    if (this.props.data) {
      data = this.props.data.device_page_resp;
      //lifecycleData = this.props.data.device_lifecycle_detail_page;
    } else {
      const detailInfo = this.state.data.toJS().detailInfo;
      data = detailInfo.data.device_page_resp;
      //lifecycleData = detailInfo.data.device_lifecycle_detail_page;
    }
    return (
      <Row>
        {renderFormDetail([
          {
            label: '固资编号',
            value: data.fixed_asset_number
          },
          {
            label: '序列号',
            value: data.sn
          },
          {
            label: '设备型号',
            value: data.model
          },
          {
            label: '设备类型',
            value: data.category
          },
          {
            label: '厂商',
            value: data.vendor
          },
          {
            label: '用途',
            value: data.usage
          },
          {
            label: '运营状态',
            value: OPERATION_STATUS[data.operation_status]
          },
          {
            label: '数据中心',
            value: data.idc && data.idc.name
          },
          {
            label: '机房管理单元',
            value: data.server_room && data.server_room.name
          },
          {
            label: '机架编号',
            value: data.server_cabinet && data.server_cabinet.number
          },
          {
            label: '机位编号',
            value: data.server_usite && data.server_usite.number
          },
          {
            label: '库房管理单元',
            value: data.store_room && data.store_room.name
          },
          {
            label: '虚拟货架',
            value: data.virtual_cabinets && data.virtual_cabinets.number
          },
          {
            label: 'RAID结构',
            value: data.raid_remark
          },
          {
            label: '硬件备注',
            value: data.hardware_remark
          },
          {
            label: '订单编号',
            value: data.order_number
          },                                    
          {
            label: '上架时间',
            value: data.onshelve_at
          },
          {
            label: '启用时间',
            value: data.started_at
          },                                              
          {
            label: '创建时间',
            value: data.created_at
          },
          {
            label: '更新时间',
            value: data.updated_at
          }
        ])}
      </Row>
    );
  }


}
