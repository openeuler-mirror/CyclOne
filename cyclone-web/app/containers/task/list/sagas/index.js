import { get, post } from 'common/xFetch2';

let created = false;
function* empty() {
  yield;
}

import { createTableSaga } from 'utils/table-saga';

const jobTableSaga = createTableSaga({
  actionNamePrefix: 'task-list/table-data',
  tableDataPath: [ 'task-list', 'tableData' ],
  datasource: '/api/cloudboot/v1/jobs',
  getExtraQuery: () => {
    return {
      rate: 'fixed_rate'
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
