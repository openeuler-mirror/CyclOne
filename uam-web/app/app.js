import Ember from 'ember';
import Resolver from 'ember/resolver';
import loadInitializers from 'ember/load-initializers';
import config from './config/environment';

let App;

Ember.MODEL_FACTORY_INJECTIONS = true;

App = Ember.Application.extend({
    modulePrefix: config.modulePrefix,
    podModulePrefix: config.podModulePrefix,
    Resolver
});

loadInitializers(App, config.modulePrefix);

// // 固定设置
// Ember.$.ajaxSetup({
//     'headers': {
//         'Accept': '*/*',
//         'Content-Type': 'application/json'
//     },
//     cache: false,
//     dataType: 'json',
//     contentType: 'application/json'
// });

export default App;
