import React from 'react';
import { Form, Button } from 'antd';
import FormElements from './FormElements';
import objectPath from './objectPath';

/**
 * [Custom]
 */
class DynamicForm extends React.Component {
  render() {
    const path = this.props.schema.name || 'root';
    return (
      <Form horizontal={true}>
        <FormElements.Object
          {...this.props.schema}
          path={path}
          key={path}
          initialValue={this.props.initialValue}
          form={this.props.form}
          updateProps={this.props.updateProps}
          dependsCallAction={this.dependsCallAction}
          onFieldChange={this.props.onFieldChange}
        />
        {this.renderSubmitButton()}
      </Form>
    );
  }

  renderSubmitButton = () => {
    return (
      <Form.Item wrapperCol={{ span: 18, offset: 4 }} className={this.props.showSubmit ? '' : 'hide'}>
        <Button
          id={`${this.props.schema.id}-submit`}
          type='primary'
          onClick={this.handleSubmit}
          style={{ float: 'right' }}
        >
          确定
        </Button>
        <Button
          id={`${this.props.schema.id}-prepare`}
          type='primary'
          onClick={this.handlePrepare}
          style={{ float: 'right' }}
          className={this.props.showPrepare ? '' : 'hide'}
        >
          预览
        </Button>
        <Button
          id={`${this.props.schema.id}-test`}
          type='primary'
          onClick={this.handleTest}
          style={{ float: 'right', marginRight: 8 }}
          className={this.props.showTest ? '' : 'hide'}
        >
          测试
        </Button>
        <Button
          id={`${this.props.schema.id}-reset`}
          type='ghost'
          onClick={this.handleReset}
          style={{ float: 'right', marginRight: 8 }}
          className={this.props.hideReset ? 'hide' : ''}
        >
          重置
        </Button>
        <Button
          id={`${this.props.schema.id}-cancel`}
          type='default'
          onClick={this.handleCancel}
          style={{ float: 'right', marginRight: 8 }}
          className={this.props.showCancel ? '' : 'hide'}
        >
          取消
        </Button>
      </Form.Item>
    );
  }

  resetFields = (names) => {
    this.props.form.resetFields(names);
  }

  onFieldChange = (ev) => {
    this.props.onFieldChange(ev);
  }

  handleReset= (e) => {
    e.preventDefault();
    this.props.form.resetFields();
  }
  handleCancel = (e) => {
    e.preventDefault();
    this.props.onCancel();
  };

  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((errors, values) => {
      if (!!errors) {
        console.log('Errors in form!!!');
        return;
      }
      this.props.onSubmit(this.revealValueObject(values[this.props.schema.name]));
    });
  }
  handleTest= (e) => {
    e.preventDefault();
    this.props.form.validateFields((errors, values) => {
      if (!!errors) {
        console.log('Errors in form!!!');
        return;
      }
      this.props.onTest(this.revealValueObject(values[this.props.schema.name]));
    });
  }
  handlePrepare = (e) => {
    e.preventDefault();
    this.props.form.validateFields((errors, values) => {
      if (!!errors) {
        console.log('Errors in form!!!');
        return;
      }
      this.props.onPrepare(this.revealValueObject(values[this.props.schema.name]));
    });
  };


  lookupDependent = (id, dependsOn) => {
    return Object.keys(dependsOn).reduce((result, key) => {
      const source = dependsOn[key];
      if (source) {
        const target = source.dependents.find((depentent) => {
          return depentent === id;
        });

        if (target && result.indexOf(key) === -1) {
          result.push(key);
        }
      }

      return result;
    }, []);
  }

  triggerDepends = (props, triggersFor) => {
    // 当前的值改变需要触发引用了该值的事件
    // 支持多依赖全部满足条件时触发事件
    const { dependsOn, form } = props;

    // 满足条件的值
    triggersFor.map((target) => {
      if (!dependsOn[target]) {
        console.error('dependsOn not found!', target);
        return;
      }
      const { dependents, action } = dependsOn[target];
      const existsValues = dependents.map((dependent) => {
        return form.getFieldValue(`${props.path}.${dependent}`);
      });
      // 如果值有undefined，说明依赖源的值不完整
      if (existsValues.indexOf(undefined) !== -1) {
        console.log('dependsOn condition not match', target, existsValues);
        return;
      }
      // call for dependsOn action
      if (action) {
        action.call(this, props, target, existsValues);
      }
    });
  }

  dependsCallAction = (props, elementId) => {
    return (e) => {
      setTimeout(() => {
        // if (e.args.length !== 1) {
        //   console.error('onFieldChange value error:', e.args);
        //   return;
        // }

        const triggersFor = this.lookupDependent(elementId, props.dependsOn);
        this.triggerDepends(props, triggersFor);
      }, 200);

      return props.onFieldChange(e);
    };
  }

  revealValueObject = (obj) => {
    const ret = {};
    const keys = Object.keys(obj);
    keys.forEach((key) => {
      if (key.indexOf('__internal__') === -1) {
        objectPath.set(ret, key, obj[key]);
      }
    });

    return ret;
  }
  // userExists(rule, value, callback) {
  //     if (!value) {
  //       callback();
  //     } else {
  //       setTimeout(() => {
  //         if (value === 'JasonWood') {
  //           callback([new Error('抱歉，该用户名已被占用。')]);
  //         } else {
  //           callback();
  //         }
  //       }, 800);
  //     }
  // }
}

export default DynamicForm;
