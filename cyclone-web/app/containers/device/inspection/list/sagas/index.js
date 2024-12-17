import { get, post } from 'common/xFetch2';

let created = false;
function* empty() {
  console.log();
}

import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'device-inspection-list/table-data',
  tableDataPath: [ 'device-inspection-list', 'tableData' ],
  datasource: '/api/cloudboot/v1/devices/inspections'
});


export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [jobTableSaga];
}
