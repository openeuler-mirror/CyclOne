import { createTableReducer } from 'utils/table-reducer';

const reducers = {
  ...createTableReducer('device/table-data', ['tableData'])
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
