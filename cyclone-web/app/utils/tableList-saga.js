import { notification } from 'antd';
import { call, put } from 'redux-saga/effects';
import { takeEvery } from 'redux-saga';
import { getWithArgs } from 'common/xFetch2';

export function createTableListSaga(options = {}) {
  const {
    actionNamePrefix = 'list',
    datasource,
    pageSize = '100'
  } = options;
  return function* defaultSaga() {
    const watches = yield [
      takeEvery(`${actionNamePrefix}/get`, getList)
    ];
  };

  function* getList(action) {
    try {
      yield put({
        type: `${actionNamePrefix}/load`
      });

      const res = yield call(getWithArgs, datasource, { page: 1, page_size: pageSize });
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }
      yield put({
        type: `${actionNamePrefix}/load/success`,
        payload: res.content.records
      });
    } catch (error) {
      console.log(error);
    }
  }
}
