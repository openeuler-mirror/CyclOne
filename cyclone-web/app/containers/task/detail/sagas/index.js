import { get, post, getWithArgs } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';

let created = false;
function* empty() {
  yield;
}


function* defaultSaga() {
  const watchers = yield [
    takeEvery('task-detail/detail-info/get', getDetailInfo)
  ];
}


function* getSettings(sns) {
  try {
    yield put({
      type: 'task-detail/devices/load'
    });

    const res = yield call(getWithArgs, `/api/cloudboot/v1/devices`, { page: 1, page_size: 100, sn: sns });
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'task-detail/devices/load/success',
      payload: res.content.records
    });
  } catch (error) {
    console.log(error);
  }
}


function* getDetailInfo(action) {
  try {
    yield put({
      type: 'task-detail/detail-info/load'
    });

    const res = yield call(get, `/api/cloudboot/v1/jobs/${action.payload.id}`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'task-detail/detail-info/load/success',
      payload: res.content
    });

    if (res.content.category === 'inspection' && JSON.stringify(res.content.target) !== '{}') {
      console.log(res.content.target);
      const obj = res.content.target;
      let sns = [];
      Object.keys(obj).forEach(key => {
        obj[key].forEach(sn => {
          sns.push(sn);
        });
      });

      yield call(getSettings, sns);
    }
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
