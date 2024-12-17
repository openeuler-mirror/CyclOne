
let created = false;
function* empty() {
  console.log();
}

import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'audit-api/table-data',
  tableDataPath: [ 'audit-api', 'tableData' ],
  datasource: '/api/cloudboot/v1/api/log'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [jobTableSaga];
}
