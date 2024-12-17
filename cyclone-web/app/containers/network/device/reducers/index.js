import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  idc: {
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
    ...createTableReducer('network-device/table-data', 'tableData'),
    ...createRegularReducer('network-device/idc', 'idc'),
    ...createRegularReducer('network-device/room', 'room')
  },
  initialState,
);

export default reducer;
