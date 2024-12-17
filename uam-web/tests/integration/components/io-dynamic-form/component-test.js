import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('io-dynamic-form', 'Integration | Component | io dynamic form', {
  integration: true
});

test('it renders', function(assert) {
  
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });" + EOL + EOL +

  this.render(hbs`{{io-dynamic-form}}`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:" + EOL +
  this.render(hbs`
    {{#io-dynamic-form}}
      template block text
    {{/io-dynamic-form}}
  `);

  assert.equal(this.$().text().trim(), 'template block text');
});
