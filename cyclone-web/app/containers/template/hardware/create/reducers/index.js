import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createTableStore, createTableReducer } from 'utils/table-reducer';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  firmwares: {
    loading: true,
    data: []
  },
  dictionaries: {
    loading: true,
    data: []
  },
  template: {
    loading: true,
    data: {
      data: []
    }
  }
});

const reducer = handleActions(
  {
    ...createRegularReducer('hardware/firmwares', 'firmwares'),
    ...createRegularReducer('hardware/dictionaries', 'dictionaries'),
    ...createRegularReducer('hardware/template', 'template'),
    'hardware/template/clear': (state, action) => {
      return state.setIn([ 'template', 'data' ], fromJS({ data: [] }))
      .setIn([ 'template', 'loading' ], true);
    },
    'hardware/template/create': (state, action) => {
      return state.setIn([ 'template', 'data' ], fromJS({ data: [] }))
        .setIn([ 'template', 'loading' ], false);
    }
  },
  initialState,
);


export default reducer;
