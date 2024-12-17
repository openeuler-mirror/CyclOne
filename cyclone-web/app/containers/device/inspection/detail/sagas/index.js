import { get, post, getWithArgs } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';

let created = false;
function* empty() {
  return;
}


function* defaultSaga() {
  const watchers = yield [
    takeEvery('inspection-detail/detail-info/get', getDetailInfo),
    takeEvery('inspection-detail/start-time/get', getStartTime)
  ];
}
function* getStartTime(action) {
  try {
    yield put({
      type: 'inspection-detail/start-time/load'
    });

    const res = yield call(get, `/api/cloudboot/v1/devices/${action.payload}/inspections/start-times`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'inspection-detail/start-time/load/success',
      payload: res.content.start_time || []
    });

  } catch (error) {
    console.log(error);
  }
}


function* getDetailInfo(action) {
  try {
    yield put({
      type: 'inspection-detail/detail-info/load'
    });

    const res = yield call(getWithArgs, `/api/cloudboot/v1/devices/${action.payload.sn}/inspections`, { start_time: action.payload.start_time });
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    console.log("see api return",res)
    yield put({
      type: 'inspection-detail/detail-info/load/success',
      //payload: res.content.result || [],
      payload: res.content || [],
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
