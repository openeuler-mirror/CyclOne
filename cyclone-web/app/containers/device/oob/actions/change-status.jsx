import React from 'react';
import { post, del, put } from 'common/xFetch2';
import { Modal, notification } from 'antd';
const confirm = Modal.confirm;

export default function action(options) {

  const typeMap = {
    powerOn: { title: '开电', api: '/api/cloudboot/v1/devices/power', method: post },
    powerOff: { title: '关电', api: '/api/cloudboot/v1/devices/power', method: del },
    reBoot: { title: '重启', api: '/api/cloudboot/v1/devices/power/restart', method: put },
    networkBoot: { title: '从网卡启动', api: '/api/cloudboot/v1/devices/power/pxe/restart', method: put },
    reAccess: { title: '重新纳管带外', api: '/api/cloudboot/v1/devices/oob/re-access', method: put }
  };

  const sns = options.records.map(it => it.sn);

  const onSubmit = () => {
    typeMap[options.type].method(typeMap[options.type].api, { sns: sns }).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      options.reload();
    });
  };

  confirm({
    title: typeMap[options.type].title,
    content: `已选择设备：${sns.length == 0 ? '全部设备' : sns.length + '台'}`,
    okText: '确定',
    cancelText: '取消',
    onOk() {
      onSubmit();
    }
  });
}
