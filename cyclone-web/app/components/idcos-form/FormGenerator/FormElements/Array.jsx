// {
//  name: 'keyArr',
//  type: 'Array',
//  itemSchema: {
//    type: 'object',
//    elements: [{
//      name: 'inputText',
//      type: 'TextInput'
//    }, {
//      name: 'inputText',
//      type: 'TextArea'
//    }]
//  }
// }

import React from 'react';
import { Form, Button, Icon } from 'antd';
import FormElements from './index';
import PropTypes from 'prop-types';
const T = PropTypes;
class ArrayType extends React.Component {

  static propTypes = {
    path: T.string.isRequired
  };
  static defaultProps = {
    initialValue: [],
    displayName: 'DynamicForm.Array'
  };


  state = {
    ...this.getInitialState()
  };

  getInitialState() {
    return this.getStateFromProps(this.props);
  }
  getStateFromProps(props) {
    const { itemSchema } = props;
    const newItemSchema = { ...itemSchema };
    const initialValue = this.props.initialValue;
    const items = [];
    for (let i = 0, l = initialValue.length; i < l; i++) {
      items.push(newItemSchema);
    }
    // 如果initialValue为空，就push一个，不然界面都看不到+号
    if (initialValue.length === 0) {
      items.push(newItemSchema);
    }
    return {
      items
    };
  }
  componentWillReceiveProps(nextProps) {
    this.setState(this.getStateFromProps(nextProps));
  }

  render() {
    const { getFieldProps, getFieldValue } = this.props.form;
    const key = this.getKey();

    const arrayField = getFieldProps(key, {
      initialValue: this.getStateFromProps(this.props).items
    });
    // title
    let $title = null;
    if (this.props.label) {
      $title = <h3>{this.props.label}</h3>;
    }
    return (
      <div className='df-array'>
        {$title}
        {this.renderElements(getFieldValue(key))}
      </div>
    );
  }

  getKey = () => {
    return '__internal__.' + this.props.path + '__arr';
  }

  renderElements = (items) => {
    const props = this.props;
    const parentPath = props.path;
    const initialValue = props.initialValue;
    return items.map((it, index) => {
      const Element = FormElements[it.type];
      if (Element) {
        const path = `${parentPath}.${index}`;
        return (
          <ArrayItem
            addItem={this.onAddItem}
            removeItem={this.onRemoveItem}
            item={it}
            itemLength={items.length}
            key={path}
          >
            <Element
              {...it}
              path={path}
              name={index + ''}
              initialValue={initialValue[index]}
              key={path}
              form={props.form}
              updateProps={props.updateProps}
              onFieldChange={props.onFieldChange}
              dependsCallAction={props.dependsCallAction}
            />
          </ArrayItem>
        );
      }
      return '';
    });
  }

  onAddItem = () => {
    const { getFieldValue, setFieldsValue } = this.props.form;
    const key = this.getKey();

    const items = getFieldValue(key);
    items.push({ ...this.props.itemSchema });

    const obj = {};
    obj[key] = items;
    setFieldsValue(obj);
  }

  onRemoveItem = (item) => {
    const { getFieldValue, setFieldsValue } = this.props.form;
    const key = this.getKey();

    let items = getFieldValue(key);
    items = items.filter(it => {
      return it !== item;
    });

    const obj = {};
    obj[key] = items;
    setFieldsValue(obj);
  }
}

class ArrayItem extends React.Component {
  render() {
    return (
      <div className='df-array-item ant-row'>
        {this.props.children}
        <span className='df-array-item-ctrl'>
          <Button type='ghost' shape='circle' onClick={this.onClickAdd}>
            <Icon type='plus' />
          </Button>
          {this.props.itemLength > 1 ? (
            <Button type='ghost' shape='circle' onClick={this.onClickMinus}>
              <Icon type='minus' />
            </Button>
          ) : (
            ''
          )}
        </span>
      </div>
    );
  }
  onClickAdd = () => {
    this.props.addItem();
  }
  onClickMinus = () => {
    this.props.removeItem(this.props.item);
  }
}

export default ArrayType;
