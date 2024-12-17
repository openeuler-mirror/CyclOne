import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';


const initialState = fromJS({
  loading: true
});

const reducer = handleActions(
  {},
  initialState,
);

export default reducer;
