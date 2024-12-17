import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  detailInfo: {
    loading: false,
    data: []
  },
  startTime: {
    loading: false,
    data: []
  }
});

const reducer = handleActions(
  {
    ...createRegularReducer('inspection-detail/detail-info', 'detailInfo'),
    ...createRegularReducer('inspection-detail/start-time', 'startTime')
  },
  initialState,
);


export default reducer;
