import React from 'react';
import FormElements from './index';

function ObjectType(props) {
  let $title = null;
  if (props.label) {
    $title = <h3>{props.label}</h3>;
  }
  return (
    <div className='df-object'>
      {$title}
      {renderElements()}
    </div>
  );

  function renderElements() {
    const elements = props.elements;
    const parentPath = props.path;
    const initialValue = props.initialValue || {};
    return elements.map(option => {
      const Element = FormElements[option.type];
      let path = `${parentPath}.${option.name}`;
      let elementInitialValue = initialValue[option.name];
      switch (option.type) {
      case 'Checkboxes':
        // checkbox is special
        path = parentPath;
        elementInitialValue = initialValue;
        break;
      case 'Array':
        // Array的表单特殊处理,需要传递这些参数
        // updateProps的更新是整个Schema，不管是否嵌套
        option.updateProps = props.updateProps;
        option.dependsCallAction = props.dependsCallAction;
        break;
      default:
        break;
      }
      if (Element) {
        return (
          <Element
            {...option}
            path={path}
            initialValue={elementInitialValue}
            key={path}
            form={props.form}
            onFieldChange={
              props.dependsOn
                ? props.dependsCallAction(props, option.id)
                : props.onFieldChange
            }
          />
        );
      }
      return null;
    });
  }
}

ObjectType.displayName = 'DynamicForm.Object';

export default ObjectType;
