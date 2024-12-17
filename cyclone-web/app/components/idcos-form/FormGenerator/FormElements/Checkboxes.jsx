import React from 'react';
import {
  Form,
  Checkbox,
  Tooltip,
  Icon
} from 'antd';
const FormItem = Form.Item;
import T from 'prop-types';
class Checkboxes extends React.Component {
  static propTypes = {
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    value: T.string,
    form: T.object.isRequired,
    path: T.string.isRequired,
    initialValue: T.object,
    disabled: T.bool,
    hidden: T.bool,
    required: T.bool,
    readonly: T.bool
  };
  static defaultProps = {
    label: '多选',
    labelColSpan: 4,
    wrapperColSpan: 18,
    initialValue: {},
    disabled: false,
    readonly: false,
    required: false,
    hidden: false
  };

  render() {
    const props = this.props;
    const itemProps = {
      label: props.label,
      labelCol: { span: props.labelColSpan },
      wrapperCol: { span: props.wrapperColSpan }
    };
    if (props.help) {
      itemProps.help = props.help;
    }
    if (props.tip) {
      itemProps.label = (
        <span>
          {props.label}
          <Tooltip title={props.tip}>
            <Icon type='question-circle-o' />
          </Tooltip>
        </span>
      );
    }
    return <FormItem {...itemProps}>{this.renderCheckboxes()}</FormItem>;
  }

  renderCheckboxes = () => {
    const props = this.props;
    const { getFieldProps } = this.props.form;
    return props.checkboxes.map(ck => {
      const fieldProps = getFieldProps(`${props.path}.${ck.name}`, {
        valuePropName: 'checked',
        rules: ck.rules,
        initialValue: props.initialValue[ck.name],
          disabled: this.props.disabled,
        onChange() {
          props.onFieldChange({
            path: `${props.path}.${ck.name}`,
            args: arguments
          });
        }
      });
      return (
        <Checkbox  disabled={this.props.disabled} {...fieldProps} className='ant-checkbox-vertical' key={ck.id}>
          {ck.label}
        </Checkbox>
      );
    });
  }
}

export default Checkboxes;
