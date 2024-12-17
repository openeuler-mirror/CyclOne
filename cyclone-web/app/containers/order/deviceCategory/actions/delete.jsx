import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { del } from 'common/xFetch2';

export default function action(options) {
  const onSubmit = () => {
    del(`/api/cloudboot/v1/device-categories`, { ids: options.records.map(it => it.id) }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
    });
  };
  confirm({
    title: `确定要删除吗?`,
    content: `共：${options.records.length} 条`,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
