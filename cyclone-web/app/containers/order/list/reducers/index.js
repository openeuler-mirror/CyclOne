import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  physicalArea: {
    data: [],
    loading: false
  },
  deviceCategory: {
    data: [],
    loading: false
  }
});

const reducer = handleActions(
  {
    ...createTableReducer('order-list/table-data', 'tableData'),
    ...createRegularReducer('order-list/physical-area', 'physicalArea'),
    ...createRegularReducer('order-list/device-categories', 'deviceCategory')
  },
  initialState,
);

export default reducer;
