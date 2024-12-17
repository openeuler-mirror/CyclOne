import React from 'react';
import {
  Form,
  Input,
  Select,
  Button,
  Checkbox,
  Radio,
  Tooltip,
  Icon
} from 'antd';
const FormItem = Form.Item;
import T from 'prop-types';
const { TextArea } = Input;

class Textarea extends React.Component {
  static propTypes={
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    hidden: T.bool,
    value: T.string,
    required: T.bool,
    disabled: T.bool,
    readonly: T.bool,
    placeholder: T.string,
    form: T.object.isRequired,
    name: T.string.isRequired,
    row: T.number,
    rules: T.array,
    path: T.string.isRequired
  };
  static defaultProps = {
    label: '文本框',
    labelColSpan: 4,
    wrapperColSpan: 18,
    disabled: false,
    readonly: false,
    required: false,
    hidden: false,
    placeholder: '请输入...',
    rows: 3
  }

  render() {
    const props = this.props;
    const { getFieldProps, getFieldError, isFieldValidating } = this.props.form;
    const fieldProps = getFieldProps(props.path, {
      rules: props.rules,
      onChange: this.onChange,
      initialValue: props.initialValue
    });
    const itemProps = {
      label: props.label,
      labelCol: { span: props.labelColSpan },
      wrapperCol: { span: props.wrapperColSpan },
      className: props.hidden ? 'hide' : ''
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
        <TextArea
          {...fieldProps}
          rows='3'
          disabled={this.props.disabled}
          placeholder={props.placeholder}
        />
      </FormItem>
    );
  }
  onChange() {
    // console.log('onChange');
  }
}

export default Textarea;
