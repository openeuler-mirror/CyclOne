import React from 'react';
import Popup from 'components/popup';
import DeviceTable from 'containers/device/common/device';

export default function action(options) {
  let query = { operation_status: 'run_with_alarm,run_without_alarm,on_shelve' };
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
        getServerRoom={true}
        limit={50}
        query={query}
        handleSubmit={(tableData) => options.handleDeviceSubmit(tableData, onSuccess)}
      />
    )
  });

}
