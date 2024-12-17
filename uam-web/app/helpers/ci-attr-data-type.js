import Ember from 'ember';

const dataTypes = {
    'string': '字符型',
    'boolean': '布尔型',
    'integer': '整型',
    'number': '浮点型',
    'enum': '枚举型',
    'datetime': '日期型',
    'date': '日期时间型'
};

export function ciAttrDataType(params) {
    let _name = dataTypes[params[0]];

    return _name ? _name : '未知数据类型';
}

export default Ember.Helper.helper(ciAttrDataType);
