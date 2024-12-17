import objectPath from 'object-path';

/**
 * [flatten description]
 * @param  {[type]} toFlatten [description]
 * @param  {[type]} prefix    [description]
 * @return {[type]}           [description]
 */
function flatten(toFlatten, prefix) {
  const result = {};
  let key;
  const traveseObject = function(theObject, path) {
    if (!(typeof theObject === 'string') && theObject instanceof Array) {
      for (let i = 0; i < theObject.length; i++) {
        // var key = path + "[" + i + "]";
        key = `${path}.${i}`;
        if (typeof theObject[i] === 'string') {
          result[key] = theObject[i];
        } else {
          traveseObject(theObject[i], key);
        }
      }
    } else {
      if (path.length > 0) {
        path = path + '.';
      }
      for (const prop in theObject) {
        if (
          theObject[prop] instanceof Object ||
          (theObject[prop] instanceof Array &&
            !(
              typeof theObject[prop] === 'string' ||
              theObject[prop] instanceof String
            ))
        ) {
          traveseObject(theObject[prop], path + prop);
        } else {
          key = path + prop;
          result[key] = theObject[prop];
        }
      }
    }
  };
  traveseObject(toFlatten, prefix);
  return result;
}
export default {
  ...objectPath,
  flatten
};
