import { getWithArgs } from 'common/xFetch2';

let created = false;
function* empty() {
  yield;
}

import { createTableListSaga } from 'utils/tableList-saga';
const roomSaga = createTableListSaga({
  actionNamePrefix: 'network-cidr/room',
  datasource: '/api/cloudboot/v1/server-rooms'
});
const networkSaga = createTableListSaga({
  actionNamePrefix: 'network-cidr/networkArea',
  datasource: '/api/cloudboot/v1/network-areas'
});
const deviceSaga = createTableListSaga({
  actionNamePrefix: 'network-cidr/device',
  datasource: '/api/cloudboot/v1/network/devices',
  pageSize: 10000
});
import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'network-cidr/table-data',
  tableDataPath: [ 'network-cidr', 'tableData' ],
  datasource: '/api/cloudboot/v1/ip-networks'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [ jobTableSaga, roomSaga, networkSaga, deviceSaga ];
}
