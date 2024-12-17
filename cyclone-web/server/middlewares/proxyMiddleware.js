/* eslint-disable global-require */
const proxy = require('http-proxy-middleware');
const logger = require('../logger');
const resolve = require('path').resolve;
const fs = require("fs");
const chalk = require('chalk');
const https = require('https');

// Proxy middleware
const addProxyMiddlewares = (app, options) => {

  let customConfig = {}, config;

  const customProxyPath = resolve(process.cwd(), '.proxy');

  if (fs.existsSync(customProxyPath)) {
    try {
      customConfig = JSON.parse(fs.readFileSync(customProxyPath, 'utf-8'));
    } catch (error) {
      logger.error("parse ./.proxy: " + error.message);
      customConfig = {};
    }
  }

  try {
    services = require(resolve(process.cwd(), 'proxy.json'));
  } catch (error) {
    logger.error("parse ./proxy.json: " + error.message);
    services = {};
  }

  try {

    const servicesKeys = Object.keys(services);
    const customKeys = Object.keys(customConfig);

    for (var attr in customConfig) {
      services[attr] = customConfig[attr];
    }

    const finalKeys = customKeys.concat(servicesKeys);
    console.log("proxy config:");
    console.log(JSON.stringify(services, 0, 4));

    finalKeys.forEach((key) => {
      const service = services[key];
      const api = service.api;
      const loglevel = service.logLevel || 'info';
      const Proxy = proxy({
        target: api,
        logLevel: loglevel,
        changeOrigin: true
        //secure: false
        // agent: new https.Agent({
        //   rejectUnauthorized: true,
        //   strictSSL: false
        // })
      });

      service.endpoints.forEach((endpoint) => {
        app.all(endpoint, function(req, res, next) {
          console.log(`${chalk.bold("->")}: ${chalk.bold(req.url)} to ${chalk.gray(api)}`);
          return Proxy(req, res, next);
        });
      });

      // log
      logger.proxyReversed(api, service.endpoints);
    });
  } catch (error) {
    logger.error("proxy config error: " + error.message);
  }
};

/**
 * Proxy middleware
 */
module.exports = (app, options) => {
  addProxyMiddlewares(app, options);
  return app;
};
