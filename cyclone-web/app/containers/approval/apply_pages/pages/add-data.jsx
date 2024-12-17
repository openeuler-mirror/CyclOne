import React from 'react';
import Popup from 'components/popup';
import ModalTable from './modal-table';

export default function action(options) {
  const onSuccess = () => {
    Popup.close();
  };
  Popup.open({
    title: `选择${options.category}`,
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <ModalTable
        {...options}
        handleSubmit={(tableData) => options.handleDeviceSubmit(tableData, onSuccess)}
      />
    )
  });

}
