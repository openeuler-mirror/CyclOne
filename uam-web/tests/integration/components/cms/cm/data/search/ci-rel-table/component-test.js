import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('cms/cm/data/search/ci-rel-table', 'Integration | Component | cms/cm/data/search/ci rel table', {
  integration: true
});

test('it renders', function(assert) {
  
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });" + EOL + EOL +

  this.render(hbs`{{cms/cm/data/search/ci-rel-table}}`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:" + EOL +
  this.render(hbs`
    {{#cms/cm/data/search/ci-rel-table}}
      template block text
    {{/cms/cm/data/search/ci-rel-table}}
  `);

  assert.equal(this.$().text().trim(), 'template block text');
});
