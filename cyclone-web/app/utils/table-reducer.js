import { fromJS } from 'immutable';

/**
 * [create create table reducer]
 * @param  {[type]} path [description]
 * @return {[type]}      [description]
 */

export function createTableStore() {
  return {
    loading: false,
    list: [],
    pagination: {
      page: 1,
      pageSize: 10
    },
    query: {},
    sorter: {},
    expand: false,
    selectedRowKeys: [],
    selectedRows: []
  };
}

function updateCache(arr = [], tableDataCache) {
  if (!(arr instanceof Array)) {
    return;
  }
  arr.forEach(it => {
    tableDataCache[(it.id === 'undefined' ? it.ID : it.id)] = it;
  });
  return tableDataCache;
}

function populateCache(ids = [], tableDataCache) {
  return ids.map(id => {
    if (tableDataCache[id]) {
      return tableDataCache[id];
    }
  });
}

/**
 * [create description]
 * @param  {[type]} path [description]
 * @return {[type]}      [description]
 */
export function createTableReducer(actionNamePrefix, keyPath = []) {
  if (typeof keyPath === 'string') {
    keyPath = [keyPath];
  }

  return {
    [actionNamePrefix + '/load']: (state, action) => {
      return state.setIn([ ...keyPath, 'loading' ], true);
    },
    [actionNamePrefix + '/load/fail']: (state, action) => {
      return state.setIn([ ...keyPath, 'loading' ], false);
    },
    [`${actionNamePrefix}/reset`]: (state, action) => {
      return state.updateIn(keyPath, tableState => {
        return fromJS(createTableStore());
      });
    },
    [`${actionNamePrefix}/toggle/expand`]: (state, action) => {
      return state.updateIn(keyPath, tableState => {
        return tableState.set('expand', !tableState.get('expand'));
      });
    },
    [`${actionNamePrefix}/load/success`]: (state, action) => {
      const payload = action.payload;
      const cache = state.getIn([ ...keyPath, 'dataCache' ]) || {};
      updateCache(payload.content || [], cache);
      return state.updateIn(keyPath, tableState => {
        return tableState
          .set('list', fromJS(payload.content))
          .set('dataCache', cache)
          .update('pagination', pagi => {
            const plainP = pagi.toJS();
            const pagiRemote = {
              ...payload.pagination
            };
            delete pagiRemote.content;

            return fromJS({
              ...plainP,
              ...pagiRemote
            });
          })
          .set('loading', false);
      });
    },
    [`${actionNamePrefix}/set/selectedRows`]: (state, action) => {
      const cache = state.getIn([ ...keyPath, 'dataCache' ]) || {};
      return state.updateIn(keyPath, tableState => {
        return (
          tableState
            // .set("selectedRows", fromJS(action.payload.selectedRows))
            .set(
              'selectedRows',
              fromJS(populateCache(action.payload.selectedRowKeys, cache))
            )
            .set('selectedRowKeys', fromJS(action.payload.selectedRowKeys))
        );
      });
    },
    [`${actionNamePrefix}/set-page-size`]: (state, action) => {
      return state.updateIn(keyPath, tableState => {
        return tableState
          .setIn([ 'pagination', 'pageSize' ], action.payload.pageSize)
          .setIn([ 'pagination', 'page' ], 1);
      });
    },
    [`${actionNamePrefix}/set-page`]: (state, action) => {
      return state.setIn(
        [ ...keyPath, 'pagination', 'page' ],
        action.payload.page
      );
    },
    [`${actionNamePrefix}/set-query`]: (state, action) => {
      return state.setIn([ ...keyPath, 'query' ], fromJS(action.payload));
    },
    [`${actionNamePrefix}/set-sorter`]: (state, action) => {
      return state.setIn([ ...keyPath, 'sorter' ], fromJS(action.payload));
    }
  };
}
