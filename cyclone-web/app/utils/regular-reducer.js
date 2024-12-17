import { fromJS } from 'immutable';

/**
 *
 * @param prefixName
 * @param keyPath
 * @returns {{}}
 *
 */
export function createRegularReducer(prefixName, keyPath = []) {
  if (typeof keyPath === 'string') {
    keyPath = [keyPath];
  }
  return {
    [`${prefixName}/load`]: (state, action) => {
      return state.setIn([ ...keyPath, 'loading' ], true);
    },
    [`${prefixName}/load/success`]: (state, action) => {
      return state.setIn([ ...keyPath, 'loading' ], false).setIn([ ...keyPath, 'data' ], fromJS(action.payload));
    }
  };
}
