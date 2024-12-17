import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';

const initialState = fromJS({
  loading: true,
  hardwareData: {},
  sysData: {},
  ip: 'no',
  inipv6: 'no',
  exipv6: 'no'
});

const reducer = handleActions(
  {
    'bunchEdit/hardware/data': (state, action) => {
      return state.set('hardwareData', fromJS(action.payload));
    },
    'bunchEdit/system/data': (state, action) => {
      return state.set('sysData', fromJS(action.payload));
    },
    'bunchEdit/ip/data': (state, action) => {
      console.log(action.payload)
      return state.set('ip', action.payload);
    },
    'bunchEdit/inipv6/data': (state, action) => {
      console.log(action.payload)
      return state.set('inipv6', action.payload);
    },
    'bunchEdit/exipv6/data': (state, action) => {
      console.log(action.payload)
      return state.set('exipv6', action.payload);
    },        
    'bunchEdit/data/clear': (state) => {
      return state.set('sysData', fromJS({})).set('hardwareData', fromJS({}));
    }
  },
  initialState,
);

export default reducer;
