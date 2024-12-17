
let created = false;
function* empty() {
  console.log();
}

import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'audit-log/table-data',
  tableDataPath: [ 'audit-log', 'tableData' ],
  datasource: '/api/cloudboot/v1/operate/log'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [jobTableSaga];
}
