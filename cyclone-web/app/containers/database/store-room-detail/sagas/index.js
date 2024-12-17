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
    takeEvery('database-store-detail/detail-info/get', getDetailInfo)
  ];
}

function* getDetailInfo(action) {
  try {
    yield put({
      type: 'database-store-detail/detail-info/load'
    });

    const res = yield call(get, `/api/cloudboot/v1/store-room/${action.payload}`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'database-store-detail/detail-info/load/success',
      payload: res.content
    });

  } catch (error) {
    console.log(error);
  }
}
import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'database-store-detail/table-data',
  tableDataPath: [ 'database-store-detail', 'tableData' ],
  datasource: '/api/cloudboot/v1/virtual-cabinets'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ defaultSaga, jobTableSaga ];
}
