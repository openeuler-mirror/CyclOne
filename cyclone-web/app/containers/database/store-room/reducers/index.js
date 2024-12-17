import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  idc: {
    loading: true,
    data: []
  }
});

const reducer = handleActions(
  {
    ...createTableReducer('database-store/table-data', 'tableData'),
    ...createRegularReducer('database-store/idc', 'idc')
  },
  initialState,
);

export default reducer;
