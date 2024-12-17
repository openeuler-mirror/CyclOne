import errorLoading from 'common/load-route-error';

async function doGet(injectReducer, injectSagas) {
  try {
    const component = await import('./index.jsx');
    return component;
  } catch (err) {
    throw err;
  }
}

export default function create(options) {
  const { errorLoading} = options;
  return function(nextState, cb) {
    doGet()
      .then(component => {
        cb(null, component.default);
      })
      .catch(errorLoading);
  };
}
