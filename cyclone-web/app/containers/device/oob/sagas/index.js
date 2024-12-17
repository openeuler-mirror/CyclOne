import { get } from 'common/xFetch2';


let created = false;
function* empty() {
  console.log();
}

import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'device-oob/table-data',
  tableDataPath: [ 'device-oob', 'tableData' ],
  datasource: '/api/cloudboot/v1/devices'
});
export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [jobTableSaga];
}
