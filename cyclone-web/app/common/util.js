/*
 * Copyright (C) 2016
 *
 * File:    util.js
 * Author:  Timothy Yeh
 * Created: 2016-10-03
 */
import { notification } from 'antd';
import { post } from 'common/xFetch2';

export function getUriParam(paramName) {
  let url = window.location.href;
  if (url.indexOf('?') != -1) {
    url = url.substr(url.indexOf('?') + 1);
    if (url.substr(url.length - 2) === '#/') {
      url = url.substr(0, url.length - 2);
    }
    const paramList = url.split('&');
    for (const param of paramList) {
      const paramMap = param.split('=');
      if (paramMap[0] === paramName) {
        return paramMap[1];
      }
    }
  }
  return null;
}

/**
 * Array2Tree
 * 转换扁平的Array为有层次关系的Tree
 *
 * @name Array2Tree
 * @function
 * @param array 扁平的数组
 * @param parent 父ID
 * @param tree 最终的树结构
 * @returns {array}
 */
export function array2Tree(array, parent, tree) {
  tree = typeof tree !== 'undefined' ? tree : [];
  parent = typeof parent !== 'undefined' ? parent : { id: '0' };

  const children = array.filter(child => {
    return child.parentId === parent.id;
  });

  if (children.length) {
    if (parent.id === '0') {
      tree = children;
    } else {
      parent.children = children;
    }

    children.forEach(child => {
      const a = array2Tree(array, child);
    });
  }

  return tree;
}

/**
 * 判断某项是否存在数组内
 *
 * @param array 数组
 * @param val 要检查的项(只支持基本)
 * @returns {boolean}
 */
export function arrayFind(array, val) {
  if (
    array === undefined ||
    array === null ||
    val === undefined ||
    val === null ||
    !(array instanceof Array)
  ) {
    return false;
  }
  for (let i = 0; i < array.length; i++) {
    if (array[i] === val) {
      return true;
    }
  }
  return false;
}

/**
 * 判断2个数组的差集
 *
 * @param array1 array1 数组
 * @returns []
 */

export function diff(array1, array2) {
  if (
    array1 === undefined ||
    array1 === null ||
    array2 === undefined ||
    array2 === null ||
    !(array1 instanceof Array) ||
    !(array2 instanceof Array)
  ) {
    return [];
  }
  return array1.concat(array2).filter(function(arg) {
    return !(array1.indexOf(arg) >= 0 && array2.indexOf(arg) >= 0);
  });
}

/**
 * lookup 获取数组的值
 *
 * @name lookup
 * @function
 * @param list
 * @param key
 * @param value
 * @param returned 返回的值字段
 * @param single 是否取只返回一条数据
 * @returns {object}
 */
export function lookup(list, key, value, returned = null, single = false) {
  if (value === undefined) {
    return;
  }

  const result = list.filter(item => {
    return item[key] === String(value);
  });

  if (!result.length) {
    return;
  }

  if (single) {
    const data = result[0];

    // 获取指定值
    if (returned) {
      return data[returned];
    }

    return data;
  }

  return result;
}

/**
 * getRandomInt
 *
 * @name getRandomInt
 * @function
 * @param min
 * @param max
 * @returns {number}
 */
export function getRandomInt(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

/**
 * numberFormat
 *
 * @name numberFormat
 * @function
 * @param number
 * @param sep
 * @returns {number}
 */
export function numberFormat(number, sep = ',') {
  const dot = String(number).split('.')[1] || '';
  return number
    .toFixed(dot.length || 2)
    .replace(/(\d)(?=(\d{3})+\.)/g, '$1' + sep);
}

export function sumBy(data, key) {
  return data.reduce((a, b) => {
    return a[key] + b[key];
  });
}

export function groupBy(data, fields, key, totalBy) {
  return data.reduce((objects, item) => {
    if (!(item[key] in objects)) {
      const obj = {};
      fields.forEach(k => {
        obj[k] = item[k];
      });
      obj[totalBy] = 1;
      objects[item[key]] = obj;
    } else {
      fields.forEach(k => {
        objects[item[key]][k] += item[k];
      });
      objects[item[key]][totalBy] += 1;
    }
    return objects;
  }, {});
}


/**
 * 校验同级目录下不能有相同的名称
 * @param data
 * @param values
 */
export const checkSameName = (data, values) => {
  if (data.children.length > 0) {
    data.children.forEach(tree => {
      if (values.pid === tree.pid && values.name === tree.name) {
        notification.error({
          message: '同级目录下不能存在同名目录'
        });
        values.flag = true;
      }
      checkSameName(tree, values);
    });
  }
  return values.flag;
};

/**
 * 删除数组中指定的项
 * @param array
 * @param item
 */
export const remove = (array, item) => {
  for (let i = 0; i < array.length; i++) {
    if (typeof item === 'string') {
      if (array[i] === item) {
        array.splice(i, 1);
      }
    } else {
      if (array[i].id === item.id) {
        array.splice(i, 1);
      }
    }

  }
  return array;
};


