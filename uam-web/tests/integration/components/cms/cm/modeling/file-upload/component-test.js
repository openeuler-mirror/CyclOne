import { moduleForComponent, test } from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('cms/cm/modeling/file-upload', 'Integration | Component | cms/cm/modeling/file upload', {
  integration: true
});

test('it renders', function(assert) {
  
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) { ... });" + EOL + EOL +

  this.render(hbs`{{cms/cm/modeling/file-upload}}`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:" + EOL +
  this.render(hbs`
    {{#cms/cm/modeling/file-upload}}
      template block text
    {{/cms/cm/modeling/file-upload}}
  `);

  assert.equal(this.$().text().trim(), 'template block text');
});
