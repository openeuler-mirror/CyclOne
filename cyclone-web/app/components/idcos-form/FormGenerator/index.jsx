import React from 'react';
import { Form } from 'antd';
import DynamicForm from './DynamicForm';
import { resolve } from 'common/resolve';
import styles from './styles.less';

const createForm = Form.create;
import T from 'prop-types';
const F = function () {};
class FormGenerator extends React.Component {
  state = {
    CustomForm: createForm()(DynamicForm)
  };
  static propTypes = {
    schema: T.shape({
      name: T.string.isRequired,
      id: T.string.isRequired,
      description: T.string,
      elements: T.array,
      onSubmit: T.func,
      onFieldChange: T.func,
      dependsOn: T.shape({
        dependents: T.array,
        triggered: T.object,
        action: T.func
      })
    }),
    showSubmit: T.bool,
    showTest: T.bool,
    showCancel: T.bool,
    showPrepare: T.bool,
    hideReset: T.bool,
    initialValue: T.object
  };
  static defaultProps = {
    schema: {
      elements: []
    },
    initialValue: {},
    showSubmit: true,
    onSubmit: F,
    onFieldChange: F
  };
  render() {
    const CustomForm = this.state.CustomForm;
    return (
      <div className={styles.formGenerator}>
        <div className='dynamic-form'>
          <CustomForm
            ref='customForm'
            {...this.props}
            onSubmit={this.onSubmit}
            onCancel={this.onCancel}
            onTest={this.onTest}
            onPrepare={this.onPrepare}
            updateProps={this.updateProps}
            onFieldChange={this.props.onFieldChange}
          />
        </div>
      </div>
    );
  }

  /**
   * [resetFields description]
   * @return {[type]} [description]
   */
  resetFields = (names) => {
    this.refs.customForm.resetFields(names);
  };

  onCancel = (values) => {
    this.props.onCancel(values);
  };
  onTest = (values) => {
    this.props.onTest(values);
  };
  onPrepare = (values) => {
    this.props.onPrepare(values);
  };

  /**
   * @param  {[type]} props  [description]
   * @param  {[type]} fields [description]
   * @return {[type]}        [description]
   */
  onFieldChange = (props, fields) => {
    this.props.onFieldChange(props, fields);
  }

  /**
   * [onSubmit description]
   * @param  {[type]} values [description]
   * @return {[type]}        [description]
   */
  onSubmit =(values) => {
    this.props.onSubmit(values);
  }

  /**
   * updateProps
   *
   * @name updateProps
   * @function
   * @param key
   * @param obj
   * @returns {undefined}
   */
  updateProps = (path, key, value) => {
    const schema = resolve(path, this.props, key, value);
    this.setState({ schema });
  }
}

export default FormGenerator;
