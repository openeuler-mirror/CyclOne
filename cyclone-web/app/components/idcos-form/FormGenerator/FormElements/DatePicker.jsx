import React from 'react';
import {
  Form,
  Tooltip,
  Icon,
  DatePicker
} from 'antd';
const FormItem = Form.Item;
import T from 'prop-types';
class FormDatePicker extends React.Component {
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
    label: '日期',
    labelColSpan: 4,
    wrapperColSpan: 18,
    disabled: false,
    readonly: false,
    required: false,
    hidden: false,
    placeholder: '请选择时间'
  };
  render() {
    const props = this.props;
    const { getFieldDecorator } = this.props.form;
    // const fieldProps = getFieldProps(props.path, {
    //   rules: props.rules,
    //   initialValue: props.initialValue
    // });

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
    return (
      <FormItem {...itemProps}>
        {getFieldDecorator(props.path, {
          rules: props.rules,
          initialValue: props.initialValue,
          disabled: this.props.disabled,
          ...spreadProps
        })(
          <DatePicker
            showTime={this.props.showTime}
            format={this.props.format || 'YYYY-MM-DD HH:mm:ss'}
            placeholder={props.placeholder}
          />
        )}
      </FormItem>
    );
  }
}

export default FormDatePicker;
