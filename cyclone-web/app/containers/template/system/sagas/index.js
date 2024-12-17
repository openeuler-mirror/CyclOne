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
    takeEvery('template-system/osFamily/get', getOsFamily)
  ];
}

import { createTableSaga } from 'utils/table-saga';

const sysTableSaga = createTableSaga({
  actionNamePrefix: 'template-system/systemConfig-table',
  tableDataPath: [ 'template-system', 'systemConfig' ],
  datasource: '/api/cloudboot/v1/system-templates'
});
const mirrorTableSaga = createTableSaga({
  actionNamePrefix: 'template-system/mirrorInstallTpl-table',
  tableDataPath: [ 'template-system', 'mirrorInstallTpl' ],
  datasource: '/api/cloudboot/v1/image-templates'
});
const networkTableSaga = createTableSaga({
  actionNamePrefix: 'template-system/network-table',
  tableDataPath: [ 'template-system', 'network' ],
  datasource: '/api/cloudboot/v1/ip-networks'
});
const oobNetworkTableSaga = createTableSaga({
  actionNamePrefix: 'template-system/oobNetwork-table',
  tableDataPath: [ 'template-system', 'oobNetwork' ],
  datasource: '/api/cloudboot/v1/oob-ip-networks'
});
function* getOsFamily(action) {
  try {
    yield put({
      type: 'template-system/osFamily/load'
    });
    const res = yield call(get, `/api/cloudboot/v1/dictionaries?type=os_family`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'template-system/osFamily/load/success',
      payload: res.content.items || []
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
  return [ sysTableSaga, mirrorTableSaga, networkTableSaga, oobNetworkTableSaga, defaultSaga ];
}
