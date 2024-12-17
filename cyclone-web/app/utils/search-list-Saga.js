import { notification } from 'antd';
import { call, put } from 'redux-saga/effects';
import { takeEvery } from 'redux-saga';
import { getWithArgs } from 'common/xFetch2';

export function createSearchListSaga(options = {}) {
  const {
    actionNamePrefix = 'list',
    datasource
  } = options;
  return takeEvery(`${actionNamePrefix}/getSearchList`, getSearchList)
  ;
  function* getSearchList(action) {
    try {
      yield put({
        type: `${actionNamePrefix}/load`
      });

      const res = yield call(getWithArgs, datasource, { page: 1, page_size: 100 });
      if (res.status !== 'success') {
        return notification.error({ message: res.message });
      }

      const list = (res.content.records || []).map(it => {
        return {
          label: it.name, value: it.id
        };
      });

      yield put({
        type: `${actionNamePrefix}/load/success`,
        payload: list
      });
    } catch (error) {
      console.log(error);
    }
  }
}
