import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';

const initialState = fromJS({
  tableData: createTableStore()
});

const reducer = handleActions(
  {
    ...createTableReducer('audit-api/table-data', 'tableData')
  },
  initialState,
);

export default reducer;
