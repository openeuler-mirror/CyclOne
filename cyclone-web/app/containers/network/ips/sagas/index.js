let created = false;
function* empty() {
  yield;
}

import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'network-ips/table-data',
  tableDataPath: [ 'network-ips', 'tableData' ],
  datasource: '/api/cloudboot/v1/ips',
  getExtraQuery: () => {
    return {
      category: 'business'
    };
  }
});

export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [jobTableSaga];
}
