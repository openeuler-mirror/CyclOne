/**
 * COMMON WEBPACK CONFIGURATION
 */

const path = require('path');
const webpack = require('webpack');

// PostCSS plugins
const cssnext = require('postcss-cssnext');
const postcssFocus = require('postcss-focus');
const postcssReporter = require('postcss-reporter');
const ProgressBarPlugin = require('progress-bar-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

const HappyPack = require('happypack');
const os = require('os');
// const happyThreadPool = HappyPack.ThreadPool({ size: 2 });
const happyThreadPool = HappyPack.ThreadPool({ size: os.cpus().length });

module.exports = options => ({
  entry: options.entry,
  output: Object.assign(
    {
      // Compile into js/build.js
      path: path.resolve(process.cwd(), 'build'),
      publicPath: '/'
    },
    options.output
  ), // Merge with env dependent settings
  module: {
    rules: [
      {
        test: /\.(js|jsx)$/,
        exclude: /node_modules/,
        // use: 'happypack/loader?id=js',
        use: [
          {
            loader: 'babel-loader',
            options: {
              presets: ['es2015-ie', 'stage-0', 'react'],
              plugins: ['transform-decorators-legacy', 'react-hot-loader/babel']
            }
          }
        ]
        // use: [
        //   {
        //     loader: "babel-loader",
        //     options: {
        //       presets: [
        //         [
        //           "es2015",
        //           {
        //             modules: false
        //           }
        //         ],
        //         "react",
        //         "stage-0"
        //       ],
        //       plugins: ["react-hot-loader/babel"]
        //     }
        //   }
        // ]
      },
      {
        test: /\.css$/,
        use: options.cssLoaders
      },
      {
        test: /\.less$/,
        use: 'happypack/loader?id=less'
        // use: [
        //   "style-loader",
        //   {
        //     loader: "css-loader",
        //     query: {
        //       importLoader: 1,
        //       localIdentName: "[path]___[name]__[local]___[hash:base64:5]",
        //       modules: true
        //     }
        //   },
        //   "postcss-loader",
        //   "less-loader"
        // ]
      },
      {
        test: /\.(jpg|png|gif|jpeg)$/,
        use: 'happypack/loader?id=pic'
        // loaders: [
        //   "file-loader?limit=5000&hash=sha512&digest=hex&size=16&name=[name].[ext]?[hash]"
        // ]
      },
      {
        test: /\.html$/,
        use: ['html-loader']
      }
    ]
  },
  plugins: options.plugins.concat([
    new HappyPack({
      id: 'js',
      verbose: true,
      threadPool: happyThreadPool,
      loaders: [
        {
          path: 'babel-loader',
          query: {
            presets: [
              [
                'es2015',
                {
                  modules: false
                }
              ],
              'react',
              'stage-0'
            ],
            plugins: ['react-hot-loader/babel']
          }
        }
      ]
    }),
    new HappyPack({
      id: 'less',
      verbose: true,
      threadPool: happyThreadPool,
      loaders: [
        'style-loader',
        {
          loader: 'css-loader',
          query: {
            importLoader: 1,
            localIdentName: '[path]___[name]__[local]___[hash:base64:5]',
            modules: true
          }
        },
        'postcss-loader',
        'less-loader'
      ]
    }),
    new HappyPack({
      id: 'pic',
      verbose: true,
      threadPool: happyThreadPool,
      loaders: [
        'file-loader?limit=5000&hash=sha512&digest=hex&size=16&name=[name].[ext]?[hash]'
      ]
    }),
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify(process.env.NODE_ENV)
      }
    }),
    new ProgressBarPlugin(),
    new CopyWebpackPlugin([
      {
        from: path.resolve(process.cwd(), 'vendor'),
        to: 'vendor'
      },
      {
        from: path.resolve(process.cwd(), 'app/assets'),
        to: 'assets'
      }
    ])
  ]),
  resolve: {
    modules: [ 'app', 'node_modules' ],
    aliasFields: ['main'],
    descriptionFiles: ['package.json'],
    mainFields: [ 'main', 'browser', 'module' ],
    extensions: [ '.js', '.jsx' ]
  },
  externals: {
    antd: 'antd',
    react: 'React',
    'react-dom': 'ReactDOM',
    'react-redux': 'ReactRedux',
    'react-router': 'ReactRouter',
    immutable: 'Immutable',
    'react-router-redux': 'ReactRouterRedux',
    'redux-saga': 'ReduxSaga',
    'redux-thunk': 'ReduxThunk',
    redux: 'Redux',
    'js-cookie': 'JsCookie',
    'object-path': 'ObjectPath',
    lodash: 'lodash',
    'redux-actions': 'ReduxActions',
    moment: 'moment',
    'react-intl': 'ReactIntl',
    reselect: 'Reselect',
    'rc-queue-anim': 'RcQueueAnim',
    echarts: 'echarts',
    'rc-animate': 'rcAnimate'
  },
  devtool: options.devtool,
  target: 'web',
  stats: {
    // Add asset Information
    assets: true,
    // Sort assets by a field
    assetsSort: 'field',
    // Add information about cached (not built) modules
    cached: true,
    // Show cached assets (setting this to `false` only shows emitted files)
    cachedAssets: true,
    // Add children information
    children: true,
    // Add chunk information (setting this to `false` allows for a less verbose output)
    chunks: true,
    // Add built modules information to chunk information
    chunkModules: true,
    // Add the origins of chunks and chunk merging info
    chunkOrigins: true,
    // Sort the chunks by a field
    chunksSort: 'field',
    // Context directory for request shortening
    // context: "../src/",
    // `webpack --colors` equivalent
    colors: true,
    // Display the distance from the entry point for each module
    depth: false,
    // Display the entry points with the corresponding bundles
    entrypoints: false,
    // Add errors
    errors: true,
    // Add details to errors (like resolving log)
    errorDetails: true,
    // Exclude modules which match one of the given strings or regular expressions
    exclude: [],
    // Add the hash of the compilation
    hash: true,
    // Set the maximum number of modules to be shown
    maxModules: 10000,
    // Add built modules information
    modules: true,
    // Sort the modules by a field
    modulesSort: 'field',
    // Show performance hint when file size exceeds `performance.maxAssetSize`
    performance: true,
    // Show the exports of the modules
    providedExports: false,
    // Add public path information
    publicPath: true,
    // Add information about the reasons why modules are included
    reasons: true,
    // Add the source code of modules
    source: true,
    // Add timing information
    timings: true,
    // Show which exports of a module are used
    usedExports: false,
    // Add webpack version information
    version: true,
    // Add warnings
    warnings: true
    // Filter warnings to be shown (since webpack 2.4.0),
    // can be a String, Regexp, a function getting the warning and returning a boolean
    // or an Array of a combination of the above. First match wins.
  }
});
