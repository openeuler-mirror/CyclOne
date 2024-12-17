import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { put } from 'common/xFetch2';

export default function action(options) {
  const records = options.records;
  const status = options.type;
  const ids = records.map(item => item.id);
  const typeMap = {
    'accepted': '验收',
    'production': '投产',
    'abolished': '裁撤'
  };
  const onSubmit = () => {
    put('/api/cloudboot/v1/idcs/status', { status: status, ids: ids }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      options.reload();
    });
  };

  confirm({
    title: `确定要${typeMap[status]}吗?`,
    content: `已选择的数据${records.length}条`,
    okText: '确定',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
