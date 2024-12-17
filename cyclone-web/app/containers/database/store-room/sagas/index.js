import { getWithArgs } from 'common/xFetch2';


let created = false;
function* empty() {
  yield;
}

import { createTableListSaga } from 'utils/tableList-saga';
const listSaga = createTableListSaga({
  actionNamePrefix: 'database-store/idc',
  datasource: '/api/cloudboot/v1/idcs'
});

import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'database-store/table-data',
  tableDataPath: [ 'database-store', 'tableData' ],
  datasource: '/api/cloudboot/v1/store-rooms'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, listSaga ];
}
