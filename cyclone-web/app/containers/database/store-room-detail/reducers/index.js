import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createRegularReducer } from 'utils/regular-reducer';
import { createTableStore, createTableReducer } from 'utils/table-reducer';


const initialState = fromJS({
  detailInfo: {
    loading: true,
    data: {}
  },
  tableData: createTableStore()
});

const reducer = handleActions(
  {
    ...createRegularReducer('database-store-detail/detail-info', 'detailInfo'),
    ...createTableReducer('database-store-detail/table-data', 'tableData')
  },
  initialState
);


export default reducer;
