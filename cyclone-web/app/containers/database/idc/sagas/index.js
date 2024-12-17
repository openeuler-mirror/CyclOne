import { get } from 'common/xFetch2';

let created = false;
function* empty() {
  console.log();
}

import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'database-idc/table-data',
  tableDataPath: [ 'database-idc', 'tableData' ],
  datasource: '/api/cloudboot/v1/idcs'
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [jobTableSaga];
}
