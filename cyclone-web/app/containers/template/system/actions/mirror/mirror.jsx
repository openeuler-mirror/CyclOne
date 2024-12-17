import React from 'react';
import Popup from 'components/popup';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import MyForm from './form';
const initDisks = [
  { "name": "/dev/sda",
    "partitions": [{ "name": "/dev/sda1", "size": "256", "fstype": "ext4", "mountpoint": "/boot" },
      { "name": "/dev/sda2", "size": "51200", "fstype": "ext4", "mountpoint": "/" },
      { "name": "/dev/sda3", "size": "2048", "fstype": "swap", "mountpoint": "swap" },
      { "name": "/dev/sda4", "size": "free", "fstype": "ext4", "mountpoint": "/home" }] },
  { "name": "/dev/sdb",
    "partitions": [{ "name": "/dev/sdb1", "size": "free", "fstype": "ext4", "mountpoint": "/disk1" }] }
];

export default function action(options) {
  const record = options.record;
  const typeMap = {
    addMirror: { title: '新增镜像配置', api: '/api/cloudboot/v1/image-templates', method: post },
    editMirror: { title: '编辑镜像配置', api: `/api/cloudboot/v1/image-templates/${record.id}`, method: put },
    copyMirror: { title: '新增镜像配置', api: '/api/cloudboot/v1/image-templates', method: post }
  };
  let initialValue = options.record;

  if (options.type === 'addMirror') {
    initialValue.disks = initDisks;
  }

  const onSubmit = (values) => {
    typeMap[options.type].method(typeMap[options.type].api, values).then(res => {
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      notification.success({ message: res.message });
      Popup.close();
      options.reload();
    });
  };

  Popup.open({
    title: `${typeMap[options.type].title}`,
    width: 800,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <MyForm
          {...options}
          initialValue={initialValue}
          id={options.record.id}
          showSubmit={true}
          initDisks={initDisks}
          onSubmit={(values) => onSubmit(values)}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
