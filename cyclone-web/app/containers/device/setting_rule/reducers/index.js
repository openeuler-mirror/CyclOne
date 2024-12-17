import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';

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
    ...createTableReducer('device-setting-rules/table-data', 'tableData'),
    'device-setting-rules/type/get': (state, action) => {
      return state.set('type', action.payload);
    }
  },
  initialState,
);

export default reducer;
