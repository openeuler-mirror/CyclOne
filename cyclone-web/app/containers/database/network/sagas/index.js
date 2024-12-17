import { getWithArgs } from 'common/xFetch2';

let created = false;
function* empty() {
  yield;
}

import { createTableListSaga } from 'utils/tableList-saga';
const listSaga = createTableListSaga({
  actionNamePrefix: 'database-network/room',
  datasource: '/api/cloudboot/v1/server-rooms'
});
import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'database-network/table-data',
  tableDataPath: [ 'database-network', 'tableData' ],
  datasource: '/api/cloudboot/v1/network-areas'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, listSaga ];
}
