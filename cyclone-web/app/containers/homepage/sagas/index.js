import { get, post, getWithArgs } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';
import { createTableSaga } from 'utils/table-saga';

let created = false;
function* empty() {
  return;
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('homepage/devices/get', getDevices),
    takeEvery('homepage/inspections/get', getInspections)
  ];
}


function* getInspections(action) {
  try {
    yield put({
      type: 'homepage/inspections/load'
    });
    const res = yield call(getWithArgs, `/api/cloudboot/v1/devices/inspections/statistics`, { ...action.payload });
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'homepage/inspections/load/success',
      payload: res.content.statistics
    });
  } catch (error) {
    console.log(error);
  }
}

function* getDevices(action) {
  try {
    yield put({
      type: 'homepage/devices/load'
    });
    const res = yield call(get, '/api/cloudboot/v1/devices/installations/statistics');
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'homepage/devices/load/success',
      payload: res.content
    });
  } catch (error) {
    console.log(error);
  }
}

/**
 * Individual exports for testing
 * @yield {[type]} [description]
 */
export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [defaultSaga];
}
