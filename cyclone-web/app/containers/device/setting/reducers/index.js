import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from "utils/regular-reducer";

const initialState = fromJS({
  tableData: createTableStore(),
  statistics: {
    loading: true,
    data: {}
  },
  type: ''
});

const reducer = handleActions(
  {
    ...createTableReducer('device-setting/table-data', 'tableData'),
    ...createRegularReducer('device-setting/statistics', 'statistics'),
    'device-setting/type/get': (state, action) => {
      return state.set('type', action.payload);
    }
  },
  initialState,
);

export default reducer;
