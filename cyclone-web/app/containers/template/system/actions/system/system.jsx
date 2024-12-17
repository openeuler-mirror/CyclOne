import React from 'react';
import Popup from 'components/popup';
import { post, put } from 'common/xFetch2';
import { notification } from 'antd';
import MyForm from './form';


export default function action(options) {
  const record = options.record;
  const typeMap = {
    addSystem: { title: '新增PXE配置', api: '/api/cloudboot/v1/system-templates', method: post },
    editSystem: { title: '编辑PXE配置', api: `/api/cloudboot/v1/system-templates/${record.id}`, method: put },
    copySystem: { title: '新增PXE配置', api: '/api/cloudboot/v1/system-templates', method: post }
  };

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
          initialValue={options.record}
          id={options.record.id}
          showSubmit={true}
          onSubmit={(values) => onSubmit(values)}
          onCancel={() => {
            Popup.close();
          }}
        />
      </div>
    )
  });
}
