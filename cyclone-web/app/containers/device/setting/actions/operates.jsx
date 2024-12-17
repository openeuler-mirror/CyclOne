import React from 'react';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;
import { put, del } from 'common/xFetch2';

export default function action(options) {
  const records = options.records;
  const typeMap = {
    reInstall: { title: '重新部署', method: put, api: '/api/cloudboot/v1/devices/installations/reinstalls', okType: 'primary', content: '' },
    cancelInstall: { title: '取消部署', method: put, api: '/api/cloudboot/v1/devices/installations/cancels', okType: 'primary', content: '' },
    finshInstall: { title: '完成部署', method: put, api: '/api/cloudboot/v1/devices/installations/setinstallsok', okType: 'primary', content: '' },
    deleteDevice: { title: '删除', method: del, api: '/api/cloudboot/v1/devices/settings', okType: 'danger', content: '，将删除设备的全部装机参数' }
  };
  const onSubmit = () => {
    const sns = records.map(it => it.sn);
    typeMap[options.type].method(typeMap[options.type].api, { sns: sns }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      options.reload();
    });
  };

  confirm({
    title: `确定要${typeMap[options.type].title}吗?`,
    content: `已选择的设备：${records.length}台${typeMap[options.type].content}`,
    okText: '确定',
    okType: `${typeMap[options.type].okType}`,
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
