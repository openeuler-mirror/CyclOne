/**
 * Create the store with asynchronously loaded reducers
 */

import { createStore, applyMiddleware, compose } from 'redux';
import { fromJS } from 'immutable';
import { routerMiddleware } from 'react-router-redux';
import createSagaMiddleware from 'redux-saga';
import createReducer from './reducers';
import thunkMiddleware from 'redux-thunk';

const sagaMiddleware = createSagaMiddleware();
const devtools = window.devToolsExtension || (() => noop => noop);

// const reducerReq = require.context('./containers', true, /reducer\.js$/);
// const sagasReq = require.context('./containers', true, /sagas\.js$/);

import globalSaga from 'containers/app/sagas';

export default function configureStore(initialState = {}, history) {
  // Create the store with two middlewares
  // 1. sagaMiddleware: Makes redux-sagas work
  // 2. routerMiddleware: Syncs the location/URL path to the state
  const middlewares = [
    thunkMiddleware,
    sagaMiddleware,
    routerMiddleware(history)
  ];

  const enhancers = [
    applyMiddleware(...middlewares),
    devtools()
  ];

  const store = createStore(
    createReducer(),
    fromJS(initialState),
    compose(...enhancers)
  );

  // Extensions
  store.runSaga = sagaMiddleware.run;
  store.asyncReducers = {}; // Async reducer registry

  store.runSaga(globalSaga);
  // // dynamic require reducers and sagas
  // sagasReq.keys().forEach((key) => {
  //   const sagas = sagasReq(key).default;
  //   sagas.map(store.runSaga);
  // });

  // reducerReq.keys().forEach((key) => {
  //   const reducer = reducerReq(key).default;
  //   const matches = key.match(/\.\/(.*)+\//);
  //   if (matches) {
  //     store.asyncReducers[reducer.NAME || matches[1]] = reducer;
  //   } else {
  //     store.asyncReducers[reducer.NAME] = reducer;
  //   }
  // });
  const nextReducers = createReducer(store.asyncReducers);
  store.replaceReducer(nextReducers);

  return store;
}
