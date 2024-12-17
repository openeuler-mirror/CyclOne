import React from 'react';
import Popup from 'components/popup';
import DeviceTable from 'containers/device/common/device';

export default function action(options) {
  const onSuccess = () => {
    Popup.close();
  };
  Popup.open({
    title: '选择设备',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <DeviceTable
        {...options}
        handleSubmit={(tableData) => options.handleDeviceSubmit(tableData, onSuccess)}
      />
    )
  });

}
