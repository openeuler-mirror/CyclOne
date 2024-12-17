import { call, put, select } from 'redux-saga/effects';
import { takeEvery } from 'redux-saga';
import { post, getWithArgs } from 'common/xFetch2';
import { createTableSaga } from 'utils/table-saga';
import * as auth from 'services/auth';
import { notification } from 'antd';


function* getUserData() {
  try {
    const req = yield call(auth.getMenuPermissions);
    const userData = req.content;
    yield put({
      type: 'global/set-user',
      payload: userData || {}
    });

    yield put({
      type: 'global/userList/get',
      payload: userData.department ? userData.department.id : ''
    });

  } catch (error) {
    console.log(error);
  }
}

function* checkLogin() {
  try {
    const token = localStorage['osinstallAuthAccessToken'];
    const hash = window.location.hash;
    if (!token && hash.indexOf('/authenticated') === -1) {
      yield put({
        type: 'global/login/first',
        payload: false
      });
      auth.ssoLogin();
    } else {
      yield put({
        type: 'global/login/first',
        payload: true
      });
    }
  } catch (error) {
    console.log(error);
  }
}

function* onLogout() {
  auth.logout();
}

function* getUserList(action) {
  try {
    yield put({
      type: 'global/userList/load'
    });

    const res = yield call(getWithArgs, `/api/cloudboot/v1/users`, { page: 1, page_size: 100, department_id: action.payload });
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    yield put({
      type: 'global/userList/load/success',
      payload: res.content.records
    });
  } catch (error) {
    console.log(error);
  }
}

import { createSearchListSaga } from 'utils/search-list-Saga';
const idcSaga = createSearchListSaga({
  actionNamePrefix: 'global/idc',
  datasource: '/api/cloudboot/v1/idcs'
});
const cabinetSaga = createSearchListSaga({
  actionNamePrefix: 'global/cabinet',
  datasource: '/api/cloudboot/v1/server-cabinets'
});
const roomSaga = createSearchListSaga({
  actionNamePrefix: 'global/room',
  datasource: '/api/cloudboot/v1/server-rooms'
});
const networkSaga = createSearchListSaga({
  actionNamePrefix: 'global/networkArea',
  datasource: '/api/cloudboot/v1/network-areas'
});
const usiteSaga = createSearchListSaga({
  actionNamePrefix: 'global/networkArea',
  datasource: '/api/cloudboot/v1/server-usites'
});
/**
 * Individual exports for testing
 * @yield {[type]} [description]
 */
function* defaultSaga() {
  const watchers = yield [
    takeEvery('global/checkLogin', checkLogin),
    takeEvery('global/getUserData', getUserData),
    takeEvery('global/logout', onLogout),
    takeEvery('global/userList/get', getUserList),
    idcSaga,
    roomSaga,
    networkSaga,
    usiteSaga,
    cabinetSaga
  ];
  return;
}

export default defaultSaga;
