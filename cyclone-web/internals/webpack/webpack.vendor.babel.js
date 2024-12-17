/**
 * WEBPACK DLL GENERATOR
 *
 * This profile is used to cache webpack's module
 * contexts for external library and framework type
 * dependencies which will usually not change often enough
 * to warrant building them from scratch every time we use
 * the webpack process.
 */

const webpack = require('webpack');
const ProgressBarPlugin = require('progress-bar-webpack-plugin');

const path = require('path');

module.exports = {
  entry: path.resolve(process.cwd(), 'app/vendor.js'),
  output: {
    filename: 'vendors.js',
    path: path.resolve(process.cwd(), 'vendor/react/'),
    library: '[name]',
    libraryTarget: 'window'
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        use: [
          {
            loader: 'babel-loader',
            options: {
              presets: ['es2015', 'es2015-ie', 'react', 'stage-0']
            }
          }
        ]
      }
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify(process.env.NODE_ENV)
      }
    }),
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        warnings: true
      }
    }),
    new ProgressBarPlugin()
  ],
  resolve: {
    modules: ['app', 'node_modules'],
    aliasFields: ['main'],
    descriptionFiles: ['package.json'],
    mainFields: ['main', 'browser', 'module'],
    extensions: ['.js', '.jsx']
  }
};
