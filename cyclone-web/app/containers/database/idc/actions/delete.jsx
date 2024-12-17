import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { del } from 'common/xFetch2';

export default function action(options) {
  const records = options.records;
  const onSubmit = () => {
    del(`/api/cloudboot/v1/idcs/${records.id}`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
    });
  };
  confirm({
    title: `确定要删除吗?`,
    content: `数据中心名称：${records.name}`,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
