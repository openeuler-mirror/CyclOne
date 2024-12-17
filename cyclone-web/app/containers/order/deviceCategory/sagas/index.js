import { get, getWithArgs } from 'common/xFetch2';
import { createTableListSaga } from 'utils/tableList-saga';

let created = false;
function* empty() {
  yield;
}


import { createTableSaga } from 'utils/table-saga';
const jobTableSaga = createTableSaga({
  actionNamePrefix: 'order-deviceCategory/table-data',
  tableDataPath: [ 'order-deviceCategory', 'tableData' ],
  datasource: '/api/cloudboot/v1/device-categories'
});


export default function create() {
  if (created) {
    return [empty];
  }
  created = true;
  return [jobTableSaga];
}
