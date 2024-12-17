import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';

const initialState = fromJS({
  tableData: createTableStore()
});

const reducer = handleActions(
  {
    ...createTableReducer('device-pre_deploy/table-data', 'tableData'),
    'device-pre_deploy/table-data/power-status/load/success': (state, action) => {
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
