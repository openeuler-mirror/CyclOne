/**
 * Created by zhangrong on 16/8/30.
 */

import xFetch from './xFetch';
import assert from './assert';

export function get(url, options) {
  return xFetch(url, options);
}

export function getWithArgs(url, args, options) {
  args = args || {};

  for (const attr in args) {
    if (args[attr] === undefined || args[attr] === null) {
      delete args[attr];
    }
  }

  let keys = Object.keys(args);
  keys = keys
    .map(key => {
      return `${key}=${args[key]}`;
    })
    .join('&');
  url = url + '?' + keys;
  return get(url, options);
}

export function post(url, data, options) {
  const opts = {
    ...options,
    method: 'POST',
    cache: 'no-cache',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  };
  return xFetch(url, opts);
}
export function postFile(url, data, options) {
  const opts = {
    ...options,
    method: 'POST',
    cache: 'no-cache',
    headers: {
      Accept: '*/*',
      'Content-Type':
        'multipart/form-data;boundary=----WebKitFormBoundaryiqw6SEM6EXa7FlBk',
      authorization: 'authorization-text'
    },
    body: JSON.stringify(data)
  };
  return xFetch(url, opts);
}

export function del(url, data, options) {
  const opts = {
    ...options,
    method: 'DELETE',
    cache: 'no-cache',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  };
  return xFetch(url, opts);
}

export function put(url, data, options) {
  const opts = {
    ...options,
    method: 'PUT',
    cache: 'no-cache',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  };
  return xFetch(url, opts);
}
