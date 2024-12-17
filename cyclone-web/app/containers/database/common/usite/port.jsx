import React from 'react';
import { Select, Col, Form, Input, Icon } from 'antd';
const FormItem = Form.Item;
import { formItemLayout } from 'common/enums';
const Option = Select.Option;

export default class Ports extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [{ name: '', port: '' }]
    };
  }


  UNSAFE_componentWillReceiveProps(props) {
    const { data } = props;
    if (data && data.length > 0) {
      this.setState({
        data: data
      });
    }
  }
  $delete = (index) => {
    const { data } = this.state;
    data.splice(index, 1);
    this.setState({ data });
  };
  add = () => {
    const { data } = this.state;
    data.push({
      name: '',
      port: ''
    });
    this.setState({ data });
  };
  render() {
    const { disabled, form, name, label } = this.props;
    const { getFieldDecorator } = form;
    const { data } = this.state;
    return (
      <div>
        {
          data.length === 0 &&
          <Icon type='plus-circle' style={{ color: 'rgb(139, 212, 109)', marginRight: 8 }} onClick={() => this.add()}/>
        }
        {
          data.length > 0 && data.map((item, index) => {
            return (
              <div>
                <Col span={12}>
                  <FormItem >
                    {getFieldDecorator(`${name}-name-${index}`, {
                      initialValue: item.name,
                      rules: [
                        {
                          required: this.props.required,
                          message: `请选择${label}名称`
                        }
                      ]
                    })(
                      <Select disabled={disabled} allowClear={true} showSearch={true}
                        filterOption={(input, option) =>
                        option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
                      }
                      >
                        {
                          (this.props.network_devices || []).map(it => <Option value={it.name}>{it.name}</Option>)
                        }
                      </Select>
                      )}
                  </FormItem>
                </Col>
                <Col span={2}>
                  <span style={{ display: 'inline-block', width: '100%', textAlign: 'center' }}>_</span>
                </Col>
                <Col span={7}>
                  <FormItem >
                    {getFieldDecorator(`${name}-port-${index}`, {
                      initialValue: item.port,
                      rules: [
                        {
                          required: this.props.required,
                          message: `请输入${label}端口`
                        }
                      ]
                    })(
                      <Input disabled={disabled} />
                      )}
                  </FormItem>
                </Col>
                <Col span={1}/>
                <Col span={2}>
                  {
                    !disabled && <div>
                      <Icon type='plus-circle' style={{ color: 'rgb(139, 212, 109)', marginRight: 8 }} onClick={() => this.add()}/>
                      {(index > 0 || !this.props.required) && <Icon type='minus-circle' style={{ color: 'rgb(255, 55, 0)' }} onClick={() => this.$delete(index)}/>}
                    </div>
                  }
                </Col>
              </div>
            );
          })
        }
      </div>
    );
  }
}
