import React from 'react';
import {
  Row
} from 'antd';
import { renderFormDetail } from 'common/utils';

export default class Detail extends React.Component {
  constructor(props) {
    super(props);
  }
  render() {
    const { loading, data } = this.props.detailInfo;
    return (
        loading ? <div>加载中...</div> :
        <Row>
          {renderFormDetail([
            {
              label: '库房名称',
              value: data.name
            },
            {
              label: '数据中心',
              value: data.idc.name
            },
            {
              label: '一级机房',
              value: data.first_server_room.name
            },
            {
              label: '城市',
              value: data.city
            },
            {
              label: '地址',
              value: data.address
            },
            {
              label: '库房负责人',
              value: data.store_room_manager
            },
            {
              label: '供应商负责人',
              value: data.vendor_manager
            }
          ])}
        </Row>

    );
  }


}
