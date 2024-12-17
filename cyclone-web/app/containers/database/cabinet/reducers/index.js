import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  network: {
    loading: true,
    data: []
  }
});

const reducer = handleActions(
  {
    ...createTableReducer('database-cabinet/table-data', 'tableData'),
    ...createRegularReducer('database-cabinet/network', 'network')
  },
  initialState,
);

export default reducer;
