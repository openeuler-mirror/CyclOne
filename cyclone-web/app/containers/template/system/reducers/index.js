import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  loading: true,
  systemConfig: createTableStore(),
  mirrorInstallTpl: createTableStore(),
  network: createTableStore(),
  oobNetwork: createTableStore(),
  osFamily: {
    data: [],
    loading: true
  }
});

const reducer = handleActions(
  {
    ...createTableReducer('template-system/systemConfig-table', 'systemConfig'),
    ...createTableReducer('template-system/mirrorInstallTpl-table', 'mirrorInstallTpl'),
    ...createTableReducer('template-system/oobNetwork-table', 'oobNetwork'),
    ...createTableReducer('template-system/network-table', 'network'),
    ...createRegularReducer('template-system/osFamily', 'osFamily')

  },
  initialState,
);

export default reducer;
