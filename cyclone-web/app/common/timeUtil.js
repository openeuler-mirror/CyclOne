/**
 * [toLocal description]
 * @param  {[type]} timestamp [description]
 * @return {[type]}           [description]
 */
import moment from 'moment';
export function dataTimeForm(data) {
  if (data) {
    let mData = moment(data);
    return mData.format('YYYY-MM-DD HH:mm:ss');
  }
}
export function dataTimeTitle(data) {
  if (data) {
    let mData = moment(data);
    return mData.format('YYYY-MM-DD-HH-mm-ss');
  }
}
export function locale(timestamp) {
  const d = new Date(timestamp);
  return d.toLocaleDateString() + ' ' + d.toLocaleTimeString();
}
