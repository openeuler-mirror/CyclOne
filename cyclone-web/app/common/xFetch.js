import 'fetch-detector';
import 'fetch-ie8';
import { base64decodeFunc } from './base64';
import * as auth from 'services/auth';
import { notification } from 'antd';

const errorMessages = res => `${res.status} ${res.statusText}`;


function check401(res) {
  if (res.status === 401) {
    notification.error({ message: '登录过期，请重新登录' });
    auth.ssoLogin();
    return Promise.reject(errorMessages(res));
  }
  return res;
}

function check404(res) {
  if (res.status === 404) {
    return Promise.reject(errorMessages(res));
  }
  return res;
}

function check500(res) {
  if (res.status === 500) {
    return Promise.reject(errorMessages(res));
  }
  return res;
}

function jsonParse(res) {
  return res.json().then(jsonResult => ({ ...res, jsonResult }));
}

function errorMessageParse(jsonResult) {
  let { status, message, content, metadata, Status } = jsonResult;
  if (status !== 'success') {
    return Promise.reject(jsonResult);
  }
  return jsonResult;
}

function xFetch(url, options) {
  const opts = { isEncode: true, ...options, credentials: 'include' };
  const locale = base64decodeFunc(localStorage.LOCALE);
  let token = localStorage.osinstallAuthAccessToken
    ? localStorage.osinstallAuthAccessToken
    : '';
  opts.headers = {
    ...opts.headers,
    'Authorization': `${token}`,
    locale
  };

  if (opts.isEncode) {
    url = encodeURI(url);
  }
  return (
    fetch(url, opts)
       .then(check401)
       //.then(check404)
      // .then(check500)
      .then(jsonParse)
      .then(res => {
        if (
          res.jsonResult &&
          res.jsonResult.statusCode &&
          res.jsonResult.statusCode.startsWith('700')
        ) {
          //检测登录超时的情况,如超时清除缓存和跳转页面
          delete localStorage.osinstallAuthAccessToken;
          auth.ssoLogin();
        }
        return res.jsonResult;
      })
  );
}

export default xFetch;
