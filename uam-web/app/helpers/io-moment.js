import Ember from 'ember';
const moment = window.moment;

export function ioMoment(params, hash) {
  var date = moment(params[0]);
  return date.format(hash.format);
}

export default Ember.Helper.helper(ioMoment);
