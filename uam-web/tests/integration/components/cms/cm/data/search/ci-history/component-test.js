import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('cms/cm/data/search/ci-history', 'Integration | Component | cms/cm/data/search/ci history', {
  integration: true
});

test('it renders', function(assert) {
  
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });" + EOL + EOL +

  this.render(hbs`{{cms/cm/data/search/ci-history}}`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:" + EOL +
  this.render(hbs`
    {{#cms/cm/data/search/ci-history}}
      template block text
    {{/cms/cm/data/search/ci-history}}
  `);

  assert.equal(this.$().text().trim(), 'template block text');
});
