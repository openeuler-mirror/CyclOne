import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  loading: false,
  userData: {},
  login: false,
  dict: {},
  userList: {
    data: [],
    loading: true
  },
  idc: {
    loading: true,
    data: []
  },
  room: {
    loading: true,
    data: []
  },
  cabinet: {
    loading: true,
    data: []
  },
  usite: {
    loading: true,
    data: []
  },
  networkArea: {
    loading: true,
    data: []
  }
});

const reducer = handleActions(
  {
    'global/set-user': (state, action) => {
      return state.set('userData', fromJS(action.payload)).set('loading', false);
    },
    'global/login/first': (state, action) => {
      return state.set('login', action.payload);
    },
    ...createRegularReducer('global/userList', 'userList'),
    ...createRegularReducer('global/idc', 'idc'),
    ...createRegularReducer('global/room', 'room'),
    ...createRegularReducer('global/cabinet', 'cabinet'),
    ...createRegularReducer('global/usite', 'usite'),
    ...createRegularReducer('global/networkArea', 'networkArea')
  },
  initialState
);

reducer.NAME = 'global';

export default reducer;
