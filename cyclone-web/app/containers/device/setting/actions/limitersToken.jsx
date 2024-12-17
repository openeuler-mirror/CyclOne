import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { del } from 'common/xFetch2';

export default function action(options) {
  const records = options.records;
  const onSubmit = () => {
    const sns = records.map(it => { return { sn: it.sn }; });
    del('/api/cloudboot/v1/devices/limiters/tokens', { tokens: sns }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      options.reload();
    });
  };

  confirm({
    title: `确定要一键释放令牌吗?`,
    content: `已选择的设备数量：${records.length}台`,
    okText: '确定',
    okType: `primary`,
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
