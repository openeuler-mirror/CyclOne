import { get } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';
import { createTableSaga } from 'utils/table-saga';


let created = false;
function* empty() {
  console.log();
}

function* defaultSaga() {
  const watchers = yield [
    takeEvery('device-special/table-data/power-status/change', getPowerStatus)
  ];
}
function* getPowerStatus(action) {
  try {
    yield put({
      type: 'device-special/table-data/power-status/load'
    });
    const res = yield call(get, `/api/cloudboot/v1/devices/${action.payload.sn}/power/status`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'device-special/table-data/power-status/load/success',
      payload: {
        sn: action.payload.sn,
        ...res.content
      }
    });
  } catch (error) {
    console.log(error);
  }
}

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'device-special/table-data',
  tableDataPath: [ 'device-special', 'tableData' ],
  datasource: '/api/cloudboot/v1/devices',
  getExtraQuery: () => {
    return {
      usage: 'SpecialDev'
    };
  }
});
export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, defaultSaga ];
}
