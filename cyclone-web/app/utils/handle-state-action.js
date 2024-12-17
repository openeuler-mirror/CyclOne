/**
 * create state action handler
 */
export default function createStateActonHandler(reducers) {
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

  return handleAction;
}
