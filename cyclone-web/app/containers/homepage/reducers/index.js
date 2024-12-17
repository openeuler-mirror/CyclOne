import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createRegularReducer } from 'utils/regular-reducer';
import { createTableStore, createTableReducer } from 'utils/table-reducer';

const initialState = fromJS({
  devices: {
    data: {},
    loading: true
  },
  inspections: {
    data: [],
    loading: true
  }
});

const reducer = handleActions(
  {
    ...createRegularReducer('homepage/inspections', 'inspections'),
    ...createRegularReducer('homepage/devices', 'devices')
  },
  initialState
);

export default reducer;
