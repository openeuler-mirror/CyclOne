/**
 * app.js
 */
import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { Router, hashHistory } from 'react-router';
import { syncHistoryWithStore } from 'react-router-redux';
import configureStore from './store';
import { translationMessages } from 'i18n';
import 'common/modal-draggable';
import 'utils/uuid';

import 'containers/app/styles/index.less';

import 'antd/lib/style/v2-compatible-reset.less';

import zhCN from 'antd/lib/locale-provider/zh_CN';


const initialState = {};
const store = configureStore(initialState, hashHistory);

if (window.devToolsExtension) {
  window.devToolsExtension.updateStore(store);
}

import { selectLocationState } from 'containers/app/selectors';
const history = syncHistoryWithStore(hashHistory, store, {
  selectLocationState: selectLocationState()
});

import App from 'containers/app';
import LanguageProvider from 'containers/common/language-provider';
import { LocaleProvider } from 'antd';

import createRoutes from './routes';
import SwalNotification from 'components/notification/SwalNotification';

window.swal = SwalNotification;
const rootRoute = {
  component: App,
  childRoutes: createRoutes(store)
};

const render = () => {
  ReactDOM.render(
    <Provider store={store}>
      <LocaleProvider locale={zhCN} messages={translationMessages}>
        <Router history={history} routes={rootRoute} />
      </LocaleProvider>
    </Provider>,
    document.getElementById('AppRoot')
  );

  setTimeout(() => {
    if (window.isIE9) {
      clearInterval(window.loadingInterval);
    }
    const $loading = document.getElementById('InitLoading');
    if ($loading) {
      $loading.style.display = 'none';
    }
  }, 10);
};

render();
