import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { post } from 'common/xFetch2';

export default function action(options) {
  const ids = options.records.map(item => item.id);
  const onSubmit = () => {
    post('/api/cloudboot/v1/server-cabinets/power', { ids: ids }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      options.reload();
    });
  };

  confirm({
    title: `确定要开电吗?`,
    content: `选择的机架个数：${ids.length}`,
    okText: '确定',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
