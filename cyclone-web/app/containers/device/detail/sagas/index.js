import { get, getWithArgs } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';

let created = false;
function* empty() {
  return;
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('device-detail/detail-info/get', getDetailInfo)
  ];
}



import { createTableListSaga } from 'utils/tableList-saga';
const idcSaga = createTableListSaga({
  actionNamePrefix: 'device-detail/idc',
  datasource: '/api/cloudboot/v1/idcs'
});


function* getDetailInfo(action) {
  try {
    yield put({
      type: 'device-detail/detail-info/load'
    });

    const res = yield call(get, `/api/cloudboot/v1/devices/${action.payload}/combined`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'device-detail/detail-info/load/success',
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
  return [ defaultSaga, idcSaga ];
}
