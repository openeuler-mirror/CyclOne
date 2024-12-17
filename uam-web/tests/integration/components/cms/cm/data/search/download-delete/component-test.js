import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('cms/cm/data/search/download-delete', 'Integration | Component | cms/cm/data/search/download delete', {
  integration: true
});

test('it renders', function(assert) {
  
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });" + EOL + EOL +

  this.render(hbs`{{cms/cm/data/search/download-delete}}`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:" + EOL +
  this.render(hbs`
    {{#cms/cm/data/search/download-delete}}
      template block text
    {{/cms/cm/data/search/download-delete}}
  `);

  assert.equal(this.$().text().trim(), 'template block text');
});
