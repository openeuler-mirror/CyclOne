import { get, post } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';

let created = false;
function* empty() {
  return;
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('hardware/firmwares/get', getFirmwares),
    takeEvery('hardware/dictionaries/get', getDictionaries),
    takeEvery('hardware/template/get', getTemplate)
  ];
}
function* getTemplate(action) {
  try {
    yield put({
      type: 'hardware/template/load'
    });

    const res = yield call(get, `/api/cloudboot/v1/hardware-templates/${action.payload}`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'hardware/template/load/success',
      payload: res.content
    });

  } catch (error) {
    console.log(error);
  }
}
function* getFirmwares(action) {
  try {
    yield put({
      type: 'hardware/firmwares/load'
    });

    const res = yield call(get, `/api/cloudboot/v1/devices/firmwares/updates/packages`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'hardware/firmwares/load/success',
      payload: res.content.packages
    });

  } catch (error) {
    console.log(error);
  }
}
function* getDictionaries(action) {
  try {
    yield put({
      type: 'hardware/dictionaries/load'
    });

    const res = yield call(get, `/api/cloudboot/v1/dictionaries?type=firmware`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'hardware/dictionaries/load/success',
      payload: res.content.items
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
  return [defaultSaga];
}
