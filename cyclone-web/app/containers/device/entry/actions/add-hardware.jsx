import React from 'react';
import Popup from 'components/popup';
import HardwareTable from './hardware';

export default function action(options) {
  const onSuccess = () => {
    Popup.close();
  };
  Popup.open({
    title: 'RAID类型',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <HardwareTable
        {...options}
        handleSubmit={(tableData) => options.handleHardwareSubmit(tableData, onSuccess)}
      />
    )
  });

}
