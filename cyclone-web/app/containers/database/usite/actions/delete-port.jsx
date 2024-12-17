import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { del } from 'common/xFetch2';

export default function action(options) {
  const records = options.records;
  const onSubmit = () => {
    del(`/api/cloudboot/v1/server-usites/${records.id}/ports`).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
    });
  };
  confirm({
    title: `确定要删除机位端口号吗?`,
    content: `机位编号：${records.number}`,
    okText: '确定',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
