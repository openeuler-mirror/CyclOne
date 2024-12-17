import React from 'react';
import {
  Form,
  InputNumber,
  Tooltip,
  Icon
} from 'antd';
const FormItem = Form.Item;
import T from 'prop-types';
class Number extends React.Component {
  static propTypes ={
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    value: T.string,
    required: T.bool,
    disabled: T.bool,
    readonly: T.bool,
    placeholder: T.string,
    form: T.object.isRequired,
    name: T.string.isRequired,
    rules: T.array,
    path: T.string.isRequired
  };
  static defaultProps = {
    label: '输入框',
    labelColSpan: 4,
    wrapperColSpan: 18,
    disabled: false,
    readonly: false,
    required: false,
    placeholder: '请输入...'
  };
  render() {
    const props = this.props;
    const { getFieldProps } = this.props.form;
    const fieldProps = getFieldProps(props.path, {
      rules: props.rules,
      onChange: this.onChange,
      initialValue: props.initialValue,
      min: props.min,
      max: props.max
    });
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
    return (
      <FormItem {...itemProps}>
        <InputNumber
          {...fieldProps}
          disabled={this.props.disabled}
          placeholder={props.placeholder}
        />
      </FormItem>
    );
  }
  onChange() {
    console.log('onChange');
  }
}

export default Number;
