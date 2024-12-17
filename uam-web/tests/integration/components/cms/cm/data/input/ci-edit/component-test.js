import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('cms/cm/data/input/ci-edit', 'Integration | Component | cms/cm/data/input/ci edit', {
  integration: true
});

test('it renders', function(assert) {
  
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });" + EOL + EOL +

  this.render(hbs`{{cms/cm/data/input/ci-edit}}`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:" + EOL +
  this.render(hbs`
    {{#cms/cm/data/input/ci-edit}}
      template block text
    {{/cms/cm/data/input/ci-edit}}
  `);

  assert.equal(this.$().text().trim(), 'template block text');
});
