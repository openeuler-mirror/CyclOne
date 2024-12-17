/**
 * get component
 */
const errorLoading = err => {
  console.error('Dynamic page loading failed', err); // eslint-disable-line no-console
};

async function doGet(injectReducer, injectSagas) {
  try {
    const saga = await import('containers/homepage/sagas/index.js');
    const reducer = await import('containers/homepage/reducers/index.js');
    injectReducer('homepage', reducer.default);
    injectSagas(saga.default());
    const component = await import('containers/homepage');
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
