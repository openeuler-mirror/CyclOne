import { createRegularReducer } from 'utils/regular-reducer';

const reducers = {
  ...createRegularReducer('device/detail', 'detailInfo')
};

const handleAction = (state, action) => {

  // action is an array
  if (action instanceof Array) {
    return action.reduce((state, action) => {
      return handleAction(state, action);
    }, state);
  }

  const type = action.type;
  if (!reducers[type]) {
    return state;
  }
  return reducers[type](state, action);
};

export default handleAction;
