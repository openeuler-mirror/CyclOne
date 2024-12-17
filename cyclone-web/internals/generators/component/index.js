/**
 * Component Generator
 */

const componentExists = require('../utils/componentExists');
const containerList = require('../utils/containerList');

module.exports = {
  description: 'Add an unconnected component',
  prompts: [{
    type: 'confirm',
    name: 'wantInContainer',
    default: false,
    message: 'Does it in container?',
  }, {
    type: 'list',
    name: 'containerName',
    message: 'Select a container',
    choices: () => containerList(),
  }, {
    type: 'list',
    name: 'type',
    message: 'Select the type of component',
    default: 'Stateless Function',
    choices: () => ['ES6 Class', 'Stateless Function'],
  }, {
    type: 'input',
    name: 'name',
    message: 'What should it be called?',
    default: 'Button',
    validate: (value) => {
      if ((/.+/).test(value)) {
        return componentExists(value) ? 'A component or container with this name already exists' : true;
      }

      return 'The name is required';
    },
  }, {
    type: 'confirm',
    name: 'wantLESS',
    default: true,
    message: 'Does it have styling?',
  }, {
    type: 'confirm',
    name: 'wantMessages',
    default: true,
    message: 'Do you want i18n messages (i.e. will this component use text)?',
  }],
  actions: (data) => {
    let name = 'app';
    if (data.wantInContainer) {
      name = `app/containers/${data.containerName}`;
    }

    // Generate index.js and index.test.js
    const actions = [{
      type: 'add',
      path: `../../${name}/components/{{properCase name}}/index.jsx`,
      templateFile: data.type === 'ES6 Class' ? './component/es6.js.hbs' : './component/stateless.js.hbs',
      abortOnFail: true,
    }, {
      type: 'add',
      path: `../../${name}/components/{{properCase name}}/tests/index.test.js`,
      templateFile: './component/test.js.hbs',
      abortOnFail: true,
    }];

    // If they want a LESS file, add styles.less
    if (data.wantLESS) {
      actions.push({
        type: 'add',
        path: `../../${name}/components/{{properCase name}}/styles.less`,
        templateFile: './component/styles.less.hbs',
        abortOnFail: true,
      });
    }

    // If they want a i18n messages file
    if (data.wantMessages) {
      actions.push({
        type: 'add',
        path: `../../${name}/components/{{properCase name}}/messages.js`,
        templateFile: './component/messages.js.hbs',
        abortOnFail: true,
      });
    }

    return actions;
  },
};
