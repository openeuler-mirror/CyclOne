import React from 'react';
import ReactDOM from 'react-dom';
import { Modal } from 'antd';
import { translationMessages } from 'i18n';
import { IntlProvider } from 'react-intl';
const div = document.createElement('div');
div.id = 'modal2';
document.body.appendChild(div);
export function open(config) {
  const El = config.content;
  const locale = config.locale || 'zh';
  ReactDOM.render(
    <IntlProvider
      locale={locale}
      key={locale}
      messages={translationMessages[locale]}
    >
      <Modal
        visible={true}
        title={config.title}
        footer={config.footer || false}
        onOk={config.onOk || null}
        width={config.width}
        height={config.height}
        maskClosable={config.maskClosable}
        wrapClassName={config.wrapClassName}
        onCancel={config.onCancel || null}
      >
        {config.content}
      </Modal>
    </IntlProvider>,
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
