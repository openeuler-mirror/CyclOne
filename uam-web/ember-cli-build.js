/*jshint node:true*/
/* global require, module */
var EmberApp = require('ember-cli/lib/broccoli/ember-app');

module.exports = function (defaults) {

    defaults.overwrite = true;

    var app = new EmberApp({
        lessOptions: {
            sourceMap: false
        },
        minifyCSS: {
            enabled: true,
            options: {
                exclude: ["vendor/**/*"]
            }
        },
        minifyJS: {
            options: {
                exclude: ["vendor/**/*"]
            }
        },
        fingerprint: {
            exclude: ["vendor/**/*"]
        },
        outputPaths: {
            app: {
                css: {
                    common: '/assets/common.css',
                    bootstrap: '/assets/bootstrap.css'
                }
            }
        },
        'ember-cli-less': {
            mergeTrees: {
                overwrite: true
            }
        }
    });

    // Use `app.import` to add additional libraries to the generated
    // output files.
    //
    // If you need to use different assets in different
    // environments, specify an object as the first parameter. That
    // object's keys should be the environment name and the values
    // should be the asset to use in that environment.
    //
    // If the library that you are including contains AMD or ES6
    // modules that you would like to import into your application
    // please specify an object with the list of modules as keys
    // along with the exports of each module as its value.

    return app.toTree();
};
