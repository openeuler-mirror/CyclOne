import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  room: {
    loading: true,
    data: []
  },
  networkArea: {
    loading: true,
    data: []
  },
  device: {
    loading: true,
    data: []
  }
});

const reducer = handleActions(
  {
    ...createTableReducer('network-cidr/table-data', 'tableData'),
    ...createRegularReducer('network-cidr/room', 'room'),
    ...createRegularReducer('network-cidr/networkArea', 'networkArea'),
    ...createRegularReducer('network-cidr/device', 'device')

  },
  initialState,
);

export default reducer;
