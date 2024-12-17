/**
 * 处理 shouldComponentUpdate 的判断逻辑。
 *
 * Created by zhangrong on 16/9/8.
 */

import { is } from 'immutable';

function shouldComponentUpdate(nextProps, nextState) {
  if (nextProps == null) {
    nextProps = {};
  }

  if (nextState == null) {
    nextState = {};
  }

  const thisProps = this.props || {};
  const thisState = this.state || {};

  if (
    Object.keys(thisProps).length !== Object.keys(nextProps).length ||
    Object.keys(thisState).length !== Object.keys(nextState).length
  ) {
    return true;
  }

  let index = 0;

  const propsKeys = Object.keys(nextProps);

  for (; index < propsKeys.length; index++) {
    const key = propsKeys[index];
    if (
      thisProps[key] !== nextProps[key] &&
      !is(thisProps[key], nextProps[key])
    ) {
      return true;
    }
  }

  const stateKeys = Object.keys(nextState);

  for (index = 0; index < stateKeys.length; index++) {
    const key = stateKeys[index];
    if (
      thisState[key] !== nextState[key] &&
      !is(thisState[key], nextState[key])
    ) {
      return true;
    }
  }

  return false;
}

export default function(component) {
  component.prototype.shouldComponentUpdate = shouldComponentUpdate;
  return component;
}
