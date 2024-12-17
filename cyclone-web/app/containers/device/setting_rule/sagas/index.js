import { get } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';


let created = false;
function* empty() {
  yield;
}

import { createTableSaga } from 'utils/table-saga';
const tableSaga = createTableSaga({
  actionNamePrefix: 'device-setting-rules/table-data',
  tableDataPath: [ 'device-setting-rules', 'tableData' ],
  datasource: '/api/cloudboot/v1/device-setting-rules',
  getExtraQuery: (state, action) => {
    const type1 = action && action.payload;
    const type = state.getIn([ 'device-setting-rules', 'type' ]) || type1;
    return {
      rule_category: type
    };
  }
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ tableSaga ];
}
