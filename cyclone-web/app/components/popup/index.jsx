import React from 'react';
import ReactDOM from 'react-dom';
import { Modal, LocaleProvider } from 'antd';
import { translationMessages } from 'i18n';
import zh_CN from 'antd/lib/locale-provider/zh_CN';

const div = document.createElement('div');
document.body.appendChild(div);
div.id = 'modal1';
export function open(config) {
  const El = config.content;
  const locale = config.locale || 'zh';
  ReactDOM.render(
    <LocaleProvider
      locale={zh_CN}
    >
      <Modal
        visible={true}
        title={config.title}
        footer={config.footer || false}
        onOk={config.onOk || null}
        width={config.width}
        height={config.height}
        maskClosable={config.maskClosable || false}
        keyboard={config.keyboard}
        wrapClassName={config.wrapClassName}
        onCancel={config.onCancel || null}
      >
        {config.content}
      </Modal>
    </LocaleProvider>,
    div
  );
}

export function close() {
  const unmountResult = ReactDOM.unmountComponentAtNode(div);
  // if (unmountResult && div.parentNode) {
  // 	div.parentNode.removeChild(div);
  // }
  // props.onCancel.apply(props, arguments);
}

export default {
  open,
  close
};
