import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  room: {
    loading: true,
    data: []
  }
});

const reducer = handleActions(
  {
    ...createTableReducer('database-network/table-data', 'tableData'),
    ...createRegularReducer('database-network/room', 'room')
  },
  initialState,
);

export default reducer;
