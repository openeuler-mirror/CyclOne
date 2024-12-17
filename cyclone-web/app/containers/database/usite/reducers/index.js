import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  cabinet: {
    loading: true,
    data: []
  },
  room: {
    loading: true,
    data: []
  }
});

const reducer = handleActions(
  {
    ...createTableReducer('database-usite/table-data', 'tableData'),
    ...createRegularReducer('database-usite/cabinet', 'cabinet'),
    ...createRegularReducer('database-usite/room', 'room')
  },
  initialState,
);

export default reducer;
