import { post, get, getWithArgs } from 'common/xFetch2';

import { createTableAsyncAction } from 'utils/table-async-action';

export default {
  ...createTableAsyncAction({
    actionNamePrefix: 'hardware/table-data',
    tableDataPath: ['tableData'],
    datasource: '/api/cloudboot/v1/hardware-templates'
  })
};
