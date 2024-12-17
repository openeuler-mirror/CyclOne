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

class TextInput extends React.Component {
  static propTypes = {
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    value: T.string,
    required: T.bool,
    disabled: T.bool,
    readonly: T.bool,
    hidden: T.bool,
    placeholder: T.string,
    form: T.object.isRequired,
    name: T.string.isRequired,
    rules: T.array,
    path: T.string.isRequired,
    inputType: T.string
  };
  static defaultProps = {
    label: '输入框',
    labelColSpan: 4,
    wrapperColSpan: 18,
    disabled: false,
    readonly: false,
    required: false,
    hidden: false,
    placeholder: '请输入...'
  };

  render() {
    const props = this.props;
    const { getFieldProps,getFieldValue } = this.props.form;
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
    const spreadProps = {};
    if (props.inputType) {
      spreadProps.type = props.inputType;
    }
    return (
      <FormItem {...itemProps}>
        <Input
          {...props}
          {...fieldProps}
          {...spreadProps}
          onBlur={()=>{
              if(typeof props.onBlur === 'function'){
                  props.onBlur(getFieldValue(props.path),this.props.form)
              }
          }
          }
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

export default TextInput;
