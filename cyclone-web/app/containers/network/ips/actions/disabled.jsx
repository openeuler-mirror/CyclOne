import React from 'react';
import { Modal, Input, notification } from 'antd';
const confirm = Modal.confirm;
import { put } from 'common/xFetch2';

export default function action(options) {
  let fixed_asset_number = options.records.fixed_asset_number;
  const setValue = (e) => {
    fixed_asset_number = e.target.value;
  };
  const onSubmit = () => {
    const ids = options.records.map(it => it.id);
    put(`/api/cloudboot/v1/ips/status/disable`, { ids: ids, remark: fixed_asset_number }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      options.reload();
    });
  };
  confirm({
    title: `确定要禁用吗?`,
    content: <div>
      <p> 已选择的IP数：{options.records.length}</p>
      <p style={{ display: 'flex', marginTop: 12, alignItems: 'baseline' }}>备注：<Input style={{ width: '80%' }} defalutValue={options.records.fixed_asset_number} onBlur={setValue}/></p>
    </div>,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
