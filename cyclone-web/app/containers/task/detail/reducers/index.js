import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  detailInfo: {
    loading: false,
    data: {}
  },
  device: {
    loading: false,
    data: []
  }
});

const reducer = handleActions(
  {
    ...createRegularReducer('task-detail/detail-info', 'detailInfo'),
    ...createRegularReducer('task-detail/devices', 'device')
  },
  initialState
);


export default reducer;
