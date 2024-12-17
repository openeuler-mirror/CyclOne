import { getWithArgs } from 'common/xFetch2';

let created = false;
function* empty() {
  yield;
}

import { createTableListSaga } from 'utils/tableList-saga';
const listSaga = createTableListSaga({
  actionNamePrefix: 'database-cabinet/network',
  datasource: '/api/cloudboot/v1/network-areas'
});

import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'database-cabinet/table-data',
  tableDataPath: [ 'database-cabinet', 'tableData' ],
  datasource: '/api/cloudboot/v1/server-cabinets'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, listSaga ];
}
