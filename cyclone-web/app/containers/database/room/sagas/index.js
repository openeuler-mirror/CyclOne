import { getWithArgs } from 'common/xFetch2';

let created = false;
function* empty() {
  console.log();
}

import { createTableListSaga } from 'utils/tableList-saga';
const listSaga = createTableListSaga({
  actionNamePrefix: 'database-room/idc',
  datasource: '/api/cloudboot/v1/idcs'
});

import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'database-room/table-data',
  tableDataPath: [ 'database-room', 'tableData' ],
  datasource: '/api/cloudboot/v1/server-rooms'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, listSaga ];
}