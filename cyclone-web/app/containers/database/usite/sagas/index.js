import { getWithArgs } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';

let created = false;
function* empty() {
  console.log();
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('database-usite/room/get', getRoom)
  ];
}

function* getRoom(action) {
  try {
    yield put({
      type: 'database-usite/room/load'
    });

    const res = yield call(getWithArgs, `/api/cloudboot/v1/server-rooms`, { page: 1, page_size: 100 });
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'database-usite/room/load/success',
      payload: res.content.records
    });
  } catch (error) {
    console.log(error);
  }
}

import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'database-usite/table-data',
  tableDataPath: [ 'database-usite', 'tableData' ],
  datasource: '/api/cloudboot/v1/server-usites'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, defaultSaga ];
}
