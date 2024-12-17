import 'babel-polyfill';
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as ReactRedux from 'react-redux';
import * as ReactRouter from 'react-router';
import * as Immutable from 'immutable';
import * as ReactRouterRedux from 'react-router-redux';
import * as ReduxSaga from 'redux-saga';
import * as ReduxThunk from 'redux-thunk';
import * as Redux from 'redux';
import * as JsCookie from 'js-cookie';
import * as ObjectPath from 'object-path';
import * as lodash from 'lodash';
import * as ReduxActions from 'redux-actions';
import * as ReactIntl from 'react-intl';
import * as Reselect from 'reselect';
import * as RcQueueAnim from 'rc-queue-anim';
import * as echarts from 'echarts';
import * as rcAnimate from 'rc-animate';
const libs = {
  React,
  ReactDOM,
  ReactRedux,
  ReactRouter,
  Immutable,
  ReactRouterRedux,
  ReduxSaga,
  ReduxThunk,
  Redux,
  JsCookie,
  ObjectPath,
  lodash,
  ReduxActions,
  ReactIntl,
  Reselect,
  RcQueueAnim,
  echarts,
  rcAnimate
};

Object.keys(libs).forEach(key => {
  window[key] = libs[key];
});
