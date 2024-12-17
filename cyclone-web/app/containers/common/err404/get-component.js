import errorLoading from 'common/load-route-error';

async function doGet() {
  try {
    const component = await import('./');
    return component;
  } catch (err) {
    throw err;
  }
}

export default function create(errorLoading, injectReducer, injectSagas) {
  return function(nextState, cb) {
    doGet()
      .then(component => {
        cb(null, component.default);
      })
      .catch(errorLoading);
  };
}
