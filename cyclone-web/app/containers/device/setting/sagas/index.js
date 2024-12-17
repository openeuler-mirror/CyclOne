import { get } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';


let created = false;
function* empty() {
  return;
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('device-setting/statistics/get', getStatistics)
  ];
}

import { createTableSaga } from 'utils/table-saga';

const tableSaga = createTableSaga({
  actionNamePrefix: 'device-setting/table-data',
  tableDataPath: [ 'device-setting', 'tableData' ],
  datasource: '/api/cloudboot/v1/devices/settings',
  getExtraQuery: (state, action) => {
    const type1 = action && action.payload;
    const type = state.getIn([ 'device-setting', 'type' ]) || type1;
    return {
      status: type
    };
  }
});

function* getStatistics() {
  try {
    yield put({
      type: 'device-setting/failure/load'
    });
    const res = yield call(get, `/api/cloudboot/v1/devices/installations/statistics`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'device-setting/statistics/load/success',
      payload: res.content
    });
  } catch (error) {
    console.log(error);
  }
}
export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ tableSaga, defaultSaga ];
}
