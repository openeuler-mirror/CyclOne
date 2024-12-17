import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';

const initialState = fromJS({
  loading: true,
  tableData: createTableStore()
});

const reducer = handleActions(
  {
    ...createTableReducer('template-hardware-list/table-data', 'tableData')
  },
  initialState,
);

export default reducer;
