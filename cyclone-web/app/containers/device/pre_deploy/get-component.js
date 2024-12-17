import errorLoading from 'common/load-route-error';

async function doGet(injectReducer, injectSagas) {
  try {
    const saga = await import('./sagas/index.js');
    const reducer = await import('./reducers/index.js');
    injectReducer('device-pre_deploy', reducer.default);
    injectSagas(saga.default());
    const component = await import('./index.jsx');
    return component;
  } catch (err) {
    throw err;
  }
}

export default function create(options) {
  const { errorLoading, injectReducer, injectSagas } = options;
  return function(nextState, cb) {
    doGet(injectReducer, injectSagas)
      .then(component => {
        cb(null, component.default);
      })
      .catch(errorLoading);
  };
}
