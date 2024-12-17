import { get } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';

let created = false;
function* empty() {
  console.log();
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('device-list/table-data/power-status/change', getPowerStatus)

  ];
}

function* getPowerStatus(action) {
  try {
    yield put({
      type: 'device-list/table-data/power-status/load'
    });
    const res = yield call(get, `/api/cloudboot/v1/devices/${action.payload.sn}/power/status`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'device-list/table-data/power-status/load/success',
      payload: {
        sn: action.payload.sn,
        ...res.content
      }
    });
  } catch (error) {
    console.log(error);
  }
}
import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'device-list/table-data',
  tableDataPath: [ 'device-list', 'tableData' ],
  datasource: '/api/cloudboot/v1/devices'
});

// import { createTableListSaga } from 'utils/tableList-saga';
// const idcSaga = createTableListSaga({
//   actionNamePrefix: 'device-list/idc',
//   datasource: '/api/cloudboot/v1/idcs'
// });
// const roomSaga = createTableListSaga({
//   actionNamePrefix: 'device-list/room',
//   datasource: '/api/cloudboot/v1/server-rooms'
// });
// const networkSaga = createTableListSaga({
//   actionNamePrefix: 'device-list/networkArea',
//   datasource: '/api/cloudboot/v1/network-areas'
// });

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, defaultSaga ];
}
