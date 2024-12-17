import React from 'react';
import { Col, Form, Slider, InputNumber, Tooltip, Icon } from 'antd';
const FormItem = Form.Item;
import T from 'prop-types';

class MSlider extends React.Component {
  static propTypes = {
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    value: T.number,
    initialValue: T.number,
    min: T.number,
    max: T.number,
    step: T.number,
    marks: T.object,
    syncInputNumber: T.bool,
    dots: T.bool,
    form: T.object.isRequired,
    path: T.string.isRequired
  };
  static defaultProps = {
    label: '选择',
    labelColSpan: 4,
    wrapperColSpan: 18,
    min: 0,
    max: 100,
    step: 1,
    dots: false,
    syncInputNumber: false
  };


  render() {
    const props = this.props;
    const { getFieldProps } = this.props.form;
    const fieldProps = getFieldProps(props.path, {
      initialValue: props.initialValue,
      rules: props.rules,
      onChange: this.onChange
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

    console.log('slider:', props, this.props.value);

    return (
      <FormItem {...itemProps}>
        <Col span={props.syncInputNumber ? 20 : 24}>
          <Slider
            {...fieldProps}
            disabled={this.props.disabled}
            min={props.min}
            max={props.max}
            step={props.step}
            marks={props.marks}
            dots={props.dots}
          />
        </Col>
        {props.syncInputNumber ? (
          <Col span={4}>
            <InputNumber
              min={props.min}
              max={props.max}
              step={props.step ? props.step : 1}
              value={this.props.value}
              onChange={this.onChange}
            />
          </Col>
        ) : (
          ''
        )}
      </FormItem>
    );
  }

  onChange = (value) => {
    this.setState({
      value
    });
  }
}

export default MSlider;
