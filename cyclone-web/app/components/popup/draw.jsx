import React from 'react';
import ReactDOM from 'react-dom';
import { Drawer, LocaleProvider } from 'antd';
import { translationMessages } from 'i18n';
import zh_CN from 'antd/lib/locale-provider/zh_CN';

const div = document.createElement('div');
document.body.appendChild(div);
div.id = 'Drawer';
export function open(config) {
  const locale = config.locale || 'zh_CN';
  ReactDOM.render(
    <LocaleProvider
      locale={zh_CN}
    >
      <Drawer
        visible={true}
        placement={config.placement || 'right'}
        destroyOnClose={true}
        title={config.title}
        width={config.width || '1280'}
        height={config.height}
        wrapClassName={config.wrapClassName}
        onClose={config.onCancel || null}
      >
        {config.content}
      </Drawer>
    </LocaleProvider>,
    div
  );
}

export function close() {
  const unmountResult = ReactDOM.unmountComponentAtNode(div);
}

export default {
  open,
  close
};
