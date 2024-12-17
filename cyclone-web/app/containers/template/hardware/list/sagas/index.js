import { get, post } from 'common/xFetch2';

let created = false;
function* empty() {
  console.log();
}

import { createTableSaga } from 'utils/table-saga';

const tableSaga = createTableSaga({
  actionNamePrefix: 'template-hardware-list/table-data',
  tableDataPath: [ 'template-hardware-list', 'tableData' ],
  datasource: '/api/cloudboot/v1/hardware-templates'
});


export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [tableSaga];
}
