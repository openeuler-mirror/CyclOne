import Ember from 'ember';

export function booleanDataType(params) {
    return params[0] === true ? '是' : '否';
}

export default Ember.Helper.helper(booleanDataType);
