import React from 'react';
import Popup from 'components/popup';
import HardwareTable from './hardware';
import SystemTable from './system';
import { Tabs, Radio, Alert, Button } from 'antd';
const TabPane = Tabs.TabPane;

export default function action(options) {

  const onSuccess = () => {
    Popup.close();
  };

  const ipChange = (e) => {
    options.dispatch({
      type: 'bunchEdit/ip/data',
      payload: e.target.value
    });
  };
  const inipv6Change = (e) => {
    options.dispatch({
      type: 'bunchEdit/inipv6/data',
      payload: e.target.value
    });
  };
  const exipv6Change = (e) => {
    options.dispatch({
      type: 'bunchEdit/exipv6/data',
      payload: e.target.value
    });
  };    
  Popup.open({
    title: '批量编辑',
    width: 1000,
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div>
        <div>
          <Alert
            message={`选择的设备： ${options.selectedRows.map(data => data.sn).join('，')}`}
            type='info'
            showIcon={true}
            style={{
              marginBottom: 10,
              overflow: 'hidden',
              wordBreak: 'break-all'
            }}
          />
        </div>
        <Tabs>
          <TabPane tab='操作系统' key='system'>
            <SystemTable
              {...options}
              bunchEdit={true}
            />
          </TabPane>
          <TabPane tab='RAID类型' key='hardware'>
            <HardwareTable
              {...options}
              bunchEdit={true}
            />
          </TabPane>
          <TabPane tab='是否分配外网IPv4' key='ip'>
            <div style={{ padding: 20 }}>
              是否分配外网IP： &nbsp;&nbsp;
              <Radio.Group onChange={(e) => ipChange(e)}>
                <Radio value='yes'>是</Radio>
                <Radio value='no'>否</Radio>
              </Radio.Group>
            </div>
          </TabPane>
          <TabPane tab='是否分配内网IPv6' key='inipv6'>
            <div style={{ padding: 20 }}>
            是否分配内网IPv6： &nbsp;&nbsp;
              <Radio.Group onChange={(e) => inipv6Change(e)}>
                <Radio value='yes'>是</Radio>
                <Radio value='no'>否</Radio>
              </Radio.Group>
            </div>
          </TabPane>
          <TabPane tab='是否分配外网IPv6' key='exipv6'>
            <div style={{ padding: 20 }}>
            是否分配外网IPv6： &nbsp;&nbsp;
              <Radio.Group onChange={(e) => exipv6Change(e)}>
                <Radio value='yes'>是</Radio>
                <Radio value='no'>否</Radio>
              </Radio.Group>
            </div>
          </TabPane>                    
        </Tabs>
        <div className='panel-footer col-right'>
          <Button type='primary' onClick={() => options.handleBunchEdit(onSuccess)}>
            确定
          </Button>
        </div>
      </div>
    )
  });

}
