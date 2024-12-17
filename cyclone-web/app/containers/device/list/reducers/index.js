import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  tableData: createTableStore(),
  // idc: {
  //   loading: true,
  //   data: []
  // },
  // room: {
  //   loading: true,
  //   data: []
  // },
  // networkArea: {
  //   loading: true,
  //   data: []
  // }
});

const reducer = handleActions(
  {
    ...createTableReducer('device-list/table-data', 'tableData'),
    // ...createRegularReducer('device-list/idc', 'idc'),
    // ...createRegularReducer('device-list/room', 'room'),
    // ...createRegularReducer('device-list/networkArea', 'networkArea'),
    'device-list/table-data/power-status/load/success': (state, action) => {
      return state.updateIn([ 'tableData', 'list' ], list => {
        let newList = list.toJS();
        newList.forEach(it => {
          if (it.sn === action.payload.sn && it.power_status !== action.payload.power_status) {
            it.power_status = action.payload.power_status;
          }
        });
        return fromJS(newList);
      });
    }
  },
  initialState,
);

export default reducer;
