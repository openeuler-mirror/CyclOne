import React from 'react';
import {
  Form,
  Radio,
  Tooltip,
  Icon
} from 'antd';
const FormItem = Form.Item;
const RadioGroup = Radio.Group;
import T from 'prop-types';

class Radios extends React.Component {
  static propTypes= {
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    value: T.string,
    rules: T.array,
    radios: T.array,
    disabled: T.bool,
    readonly: T.bool,
    hidden: T.bool,
    required: T.bool,
    form: T.object.isRequired,
    path: T.string.isRequired
  };
  static defaultProps = {
    label: '单选',
    labelColSpan: 4,
    wrapperColSpan: 18,
    disabled: false,
    readonly: false,
    required: false,
    hidden: false,
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
    return <FormItem {...itemProps}>{this.renderRadios()}</FormItem>;
  }

  renderRadios() {
    const props = this.props;
    const { getFieldProps } = this.props.form;
    const fieldProps = getFieldProps(props.path, {
      rules: props.rules,
      initialValue: props.initialValue,
      disabled: this.props.disabled,
    });
    if (props.defaultValue) {
      fieldProps.defaultValue = props.defaultValue;
    }

    return (
      <RadioGroup {...fieldProps}>
        {props.radios.map(r => {
          return (
            <Radio disabled={this.props.disabled} value={r.value} key={r.value + r.label + props.id}>
              {r.label}
            </Radio>
          );
        })}
      </RadioGroup>
    );
  }
}

export default Radios;
