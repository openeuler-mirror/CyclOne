import React from 'react';
import Popup from 'components/popup';
import SystemTable from './system';

export default function action(options) {
  const onSuccess = () => {
    Popup.close();
  };
  Popup.open({
    title: '操作系统',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <SystemTable
        {...options}
        handleSubmit={(tableData) => options.handleSysTemplateSubmit(tableData, onSuccess)}
      />
    )
  });

}
