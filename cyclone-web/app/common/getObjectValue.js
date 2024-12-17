/**
 * [description]
 * @return {[type]} [description]
 */
export function getObjectValue(obj, array) {
  if (obj == null) {
    return undefined;
  }
  let len;
  if (array == undefined) {
    len = 0;
  } else {
    len = array.length;
  }
  let res = obj;
  for (let i = 0; i < len; i++) {
    res = res[array[i]];
    if (res == null) {
      return undefined;
    }
  }
  return res;
}

export function getArrayLength(obj, array) {
  var arr = getObjectValue(obj, array);
  if (arr == undefined) {
    return 0;
  } else {
    return arr.length;
  }
}
