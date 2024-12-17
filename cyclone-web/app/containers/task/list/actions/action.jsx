import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { put, del } from 'common/xFetch2';

export default function action(options) {
  const job_id = options.record.id;
  const type = options.type;
  const typeMap = {
    '_pause': { api: `/api/cloudboot/v1/jobs/${job_id}/pausing`, title: '暂停', method: put, okType: 'primary' },
    '_continue': { api: `/api/cloudboot/v1/jobs/${job_id}/unpausing`, title: '继续', method: put, okType: 'primary' },
    '_delete': { api: `/api/cloudboot/v1/jobs/${job_id}`, title: '继续', method: del, okType: 'danger' }
  };
  const onSubmit = () => {
    typeMap[type].method(typeMap[type].api).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: '操作成功' });
      options.reload();
    });
  };

  confirm({
    title: `确定要${typeMap[type].title}该任务吗?`,
    content: `任务标题：${options.record.title}`,
    okText: '确定',
    okType: typeMap[type].okType,
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
