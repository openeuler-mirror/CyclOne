import { get, getWithArgs } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';
import { createTableListSaga } from 'utils/tableList-saga';

let created = false;
function* empty() {
  yield;
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('order-list/physical-area/get', getPhysicalArea)
  ];
}


import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'order-list/table-data',
  tableDataPath: ['order-list', 'tableData'],
  datasource: '/api/cloudboot/v1/orders'
});

const deviceSaga = createTableListSaga({
  actionNamePrefix: 'order-list/device-categories',
  datasource: '/api/cloudboot/v1/device-categories'
});

function* getPhysicalArea() {
  try {
    yield put({
      type: 'order-list/physical-area/load'
    });
    const res = yield call(get, `/api/cloudboot/v1/devices/query-params/physical_area`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'order-list/physical-area/load/success',
      payload: res.content.list
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
  return [defaultSaga, jobTableSaga, deviceSaga];
}
