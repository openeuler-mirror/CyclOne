import { post, get, getWithArgs } from 'common/xFetch2';

import { createTableAsyncAction } from 'utils/table-async-action';

export default {
  ...createTableAsyncAction({
    actionNamePrefix: 'system/table-data',
    tableDataPath: ['tableData'],
    noMoreQuery: true,
    datasource: '/api/cloudboot/v1/os-templates',
    getExtraQuery: () => {
      return {
        os_lifecycle: 'active_default,active'
      };
    }
  })
};
