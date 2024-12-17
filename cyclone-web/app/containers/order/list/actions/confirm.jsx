import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { put } from 'common/xFetch2';

export default function action(options) {
  const onSubmit = () => {
    put(`/api/cloudboot/v1/order/status`, { "id": options.records.id, "status": "finished" }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
    });
  };
  confirm({
    title: `确定订单已完成?`,
    content: `订单号：${options.records.number}`,
    okText: '确定',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
