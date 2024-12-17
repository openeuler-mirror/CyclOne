import React from 'react';
import { Form, Switch, Tooltip, Icon } from 'antd';
const FormItem = Form.Item;
import T from 'prop-types';

class MSwitch extends React.Component {
  static propTypes= {
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    value: T.bool,
    initialValue: T.bool,
    checkedChildren: T.string,
    unCheckedChildren: T.string,
    size: T.string,
    form: T.object.isRequired,
    path: T.string.isRequired
  };
  static defaultProps = {
    label: '选择',
    labelColSpan: 4,
    wrapperColSpan: 18,
    initialValue: false
  };
  render() {
    const props = this.props;
    const { getFieldProps } = this.props.form;
    const fieldProps = getFieldProps(props.path, {
      initialValue: props.initialValue,
      rules: props.rules,
      onChange() {
        props.onFieldChange({
          path: `${props.path}.${props.name}`,
          args: arguments
        });
      }
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
        <Switch
          {...fieldProps}
          size={props.size}
          defaultChecked={props.initialValue}
          checkedChildren={props.checkedChildren}
          unCheckedChildren={props.unCheckedChildren}
          disabled={this.props.disabled}
        />
      </FormItem>
    );
  }
}

export default MSwitch;
