/* eslint consistent-return:0 */
const express = require('express');
const logger = require('./logger');
const argv = require('minimist')(process.argv.slice(2));
const resolve = require('path').resolve;
const app = express();
const bodyParser = require('body-parser');
// require('dotenv').config();

const FrontendServer = function(app) {
  const frontend = require('./middlewares/frontendMiddleware');
  const proxy = require('./middlewares/proxyMiddleware');
  // If you need a backend, e.g. an API, add your custom backend-specific middleware here
  // app.use('/api', myApi);

  // Proxy back-end API
  proxy(app, {
    config: resolve(process.cwd(), 'proxy.json')
  });

  // In production we need to pass these values in instead of relying on webpack
  frontend(app, {
    outputPath: resolve(process.cwd(), 'build'),
    publicPath: '/'
  });
};

const MockServer = function(app) {
  const globSync = require('glob').sync;
  const mocks = globSync('./mocks/**/*.js', {
    cwd: __dirname
  }).map(require);
  const morgan = require('morgan');

  app.use(morgan('dev'));
  mocks.forEach(route => {
    route(app);
  });
};

const Run = function() {
  const isDev = process.env.NODE_ENV !== 'production';
  const ngrok =
    (isDev && process.env.ENABLE_TUNNEL) || argv.tunnel
      ? require('ngrok')
      : false;
  // get the intended port number, use port 3000 if not provided
  const port = argv.port || process.env.PORT || 3000;

  if (process.env.MOCK === 'true') {
    MockServer(app);
  }

  FrontendServer(app);

  // Start your app.
  app.listen(port, err => {
    if (err) {
      return logger.error(err.message);
    }

    // Connect to ngrok in dev mode
    if (ngrok) {
      ngrok.connect(
        port,
        (innerErr, url) => {
          if (innerErr) {
            return logger.error(innerErr);
          }

          logger.appStarted(port, url);
        }
      );
    } else {
      logger.appStarted(port);
    }
  });
};

Run();
