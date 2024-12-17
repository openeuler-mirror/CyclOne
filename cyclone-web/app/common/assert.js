/*
 * Copyright (C) 2016
 *
 * File:    assert.js
 * Author:  Timothy Yeh
 * Created: 2016-10-03
 *
 * 断言函数
 */

export default function assert(condition, message) {
  if (!condition) {
    message = message || 'Assertion failed';
    if (typeof Error !== 'undefined') {
      throw new Error(message);
    }
    throw message; // Fallback
  }
}
