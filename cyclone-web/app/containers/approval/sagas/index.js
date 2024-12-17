import { get, getWithArgs } from 'common/xFetch2';
import { takeEvery } from 'redux-saga';
import { call, put } from 'redux-saga/effects';
import { notification } from 'antd';

let created = false;
function* empty() {
  console.log();
}

function* defaultSaga() {
  yield;
}


import { createTableSaga } from 'utils/table-saga';

const pendingTableSaga = createTableSaga({
  actionNamePrefix: 'approval/pending-table-data',
  tableDataPath: [ 'approval', 'pendingTableData' ],
  datasource: '/api/cloudboot/v1/users/pending/approvals'
});

const approvedTableSaga = createTableSaga({
  actionNamePrefix: 'approval/approved-table-data',
  tableDataPath: [ 'approval', 'approvedTableData' ],
  datasource: '/api/cloudboot/v1/users/approved/approvals'
});

const initiatedTableSaga = createTableSaga({
  actionNamePrefix: 'approval/initiated-table-data',
  tableDataPath: [ 'approval', 'initiatedTableData' ],
  datasource: '/api/cloudboot/v1/users/initiated/approvals'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ defaultSaga, pendingTableSaga, approvedTableSaga, initiatedTableSaga ];
}
