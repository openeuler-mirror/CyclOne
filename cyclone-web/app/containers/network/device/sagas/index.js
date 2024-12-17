import { getWithArgs } from 'common/xFetch2';

let created = false;
function* empty() {
  yield;
}
import { createTableListSaga } from 'utils/tableList-saga';
const idcSaga = createTableListSaga({
  actionNamePrefix: 'network-device/idc',
  datasource: '/api/cloudboot/v1/idcs'
});
const roomSaga = createTableListSaga({
  actionNamePrefix: 'network-device/room',
  datasource: '/api/cloudboot/v1/server-rooms'
});


import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'network-device/table-data',
  tableDataPath: [ 'network-device', 'tableData' ],
  datasource: '/api/cloudboot/v1/network/devices'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, idcSaga, roomSaga ];
}
