import React from 'react';
import { Input, Switch, Select, Col, Button, Form } from 'antd';
const Option = Select.Option;

const FormItem = Form.Item;
const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 8 }
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 16 }
  }
};
class Template extends React.Component {


  render() {
    const { getFieldDecorator, getFieldValue } = this.props.form;
    const { data, index, disabled } = this.props;
    const actions = data.action || [];
    const metaData = data.metaData || {};

    let custom = 'custom';
    let bios_key = {};
    let action = '';

    //raid,oob action选项联动
    if (data.category === 'raid' || data.category === 'oob') {
      action = getFieldValue(`${data.category}-action-${data.uuid}`) || data.initialValue.action;
    }

    //bios,key选项联动
    if (data.category === 'bios') {
      custom = getFieldValue(`${data.category}-custom-${data.uuid}`) || data.initialValue.metadata.manufacturer || 'custom';
      const key = getFieldValue(`${data.category}-key-${data.uuid}`);
      if (custom && custom !== 'custom') {
        metaData[custom].map(mapKey => {
          if (mapKey.ci === key) {
            bios_key = mapKey;
          }
        });
      }
    }

    return (
      <div>
        <div className='configure-content'>
          <Col span='2'>{data.category.toLocaleUpperCase()}</Col>

          {/*firmware类型时第二列不显示固定文字*/}
          {
            data.category !== 'firmware' &&
            <Col span='3'>
              {typeof data.block === 'string' ? data.block :
              <FormItem>
                {getFieldDecorator(`${data.category}-custom-${data.uuid}`, {
                  initialValue: data.initialValue.metadata.manufacturer || 'custom'
                })(
                  <Select disabled={disabled} style={{ width: '100%' }}>
                    {data.block.options.map((item, index) => <Option key={index} value={item.value}>{item.name}</Option>)}
                  </Select>
                    )}
              </FormItem>
              }
            </Col>
          }

          {/*bios类型*/}
          {
            data.category === 'bios' && (custom === 'custom' ?
              (<div><Col span='8'>
                <FormItem label='参数名' {...formItemLayout}>
                  {getFieldDecorator(`${data.category}-key-${data.uuid}`, {
                    initialValue: data.initialValue.metadata.key,
                    rules: [{
                      required: true,
                      message: 'key不能为空'
                    }]
                  })(
                    <Input disabled={disabled} />
                  )}
                </FormItem>
              </Col>
                <Col span='8'>
                  <FormItem label='参数值' {...formItemLayout}>
                    {getFieldDecorator(`${data.category}-value-${data.uuid}`, {
                      initialValue: data.initialValue.metadata.value,
                      rules: [{
                        required: true,
                        message: 'value不能为空'
                      }]
                    })(
                      <Input disabled={disabled} />
                  )}
                  </FormItem>
                </Col>
              </div>
              ) : (
                <div><Col span='8'>
                  <FormItem label='参数名' {...formItemLayout}>
                    {getFieldDecorator(`${data.category}-key-${data.uuid}`, {
                      initialValue: data.initialValue.metadata.key,
                      rules: [{
                        required: true,
                        message: 'key不能为空'
                      }]
                    })(
                      <Select disabled={disabled} style={{ width: '100%' }}>
                        {(metaData[custom] || []).map((item, index) => <Option key={index} value={item.ci}>{item.name}</Option>)}
                      </Select>
                    )}
                  </FormItem>
                </Col>
                  <Col span='8'>
                    {
                    bios_key.type === 'radio' &&
                    <FormItem label='参数值' {...formItemLayout}>
                      {getFieldDecorator(`${data.category}-value-${data.uuid}`, {
                        initialValue: data.initialValue.metadata.value === 'enable' || data.initialValue.metadata.value === 'Enabled',
                        valuePropName: 'checked',
                        rules: [{
                          required: true,
                          message: 'value不能为空'
                        }]
                      })(
                        <Switch disabled={disabled} checkedChildren={custom === 'Dell' ? 'enable' : 'Enabled'} unCheckedChildren={custom === 'Dell' ? 'disable' : 'Disabled'} />
                      )}
                    </FormItem>
                  }
                    {
                    bios_key.type === 'select' &&
                    <FormItem label='参数值' {...formItemLayout}>
                      {getFieldDecorator(`${data.category}-value-${data.uuid}`, {
                        initialValue: data.initialValue.metadata.value,
                        rules: [{
                          required: true,
                          message: 'value不能为空'
                        }]
                      })(
                        <Select disabled={disabled} style={{ width: '100%' }}>
                          {(bios_key.options || []).map((item, index) => <Option key={index} value={item.id}>{item.name}</Option>)}
                        </Select>
                      )}
                    </FormItem>
                  }
                  </Col>
                </div>
              ))
          }

          {/*firmware类型*/}
          { data.category === 'firmware' &&
          (<div>
            <Col span='6'>
              <FormItem label='待固件升级包' {...formItemLayout}>
                {getFieldDecorator(`${data.category}-file-${data.uuid}`, {
                  initialValue: data.initialValue.metadata.file,
                  rules: [{
                    required: true,
                    message: '待固件升级包不能为空'
                  }]
                })(
                  <Select disabled={disabled} style={{ width: '100%' }} placeholder={'请选择'} >
                    {(this.props.firmwares || []).map((item, index) => <Option key={index} value={item}>{item}</Option>)}
                  </Select>
                )}
              </FormItem>
            </Col>
            <Col span='6'>
              <FormItem label='固件类型' {...formItemLayout}>
                {getFieldDecorator(`${data.category}-category-${data.uuid}`, {
                  initialValue: data.initialValue.metadata.category ? JSON.stringify({ name: data.initialValue.metadata.category, value: data.initialValue.metadata.category_desc }) : '',
                  rules: [{
                    required: true,
                    message: '固件类型不能为空'
                  }]
                })(
                  <Select disabled={disabled} style={{ width: '100%' }} placeholder={'请选择'} >
                    {(this.props.dictionaries || []).map((item, index) => <Option key={index} value={JSON.stringify({ name: item.name, value: item.value })}>{item.name}</Option>)}
                  </Select>
                )}
              </FormItem>
            </Col>
            <Col span='6'>
              <FormItem label='固件版本' {...formItemLayout}>
                {getFieldDecorator(`${data.category}-expected-${data.uuid}`, {
                  initialValue: data.initialValue.metadata.expected,
                  rules: [{
                    required: true,
                    message: '固件版本不能为空'
                  }]
                })(
                  <Input disabled={disabled} />
                )}
              </FormItem>
            </Col>
          </div>
          )
          }

          {/*raid 和 oob 类型显示actions列表*/}
          { (data.category === 'raid' || data.category === 'oob') &&
          <Col span='6'>
            <FormItem>
              {getFieldDecorator(`${data.category}-action-${data.uuid}`, {
                initialValue: data.initialValue.action,
                rules: [{
                  required: true,
                  message: ' '
                }]
              })(
                <Select disabled={disabled} style={{ width: '100%' }} placeholder={'请选择'} >
                  {actions.map((item, index) => <Option key={index} value={item.value}>{item.name}</Option>)}
                </Select>
              )}
            </FormItem>
          </Col>
          }

          {/*raid 和 oob 类型显示*/}
          {
            (data.category === 'raid' || data.category === 'oob') && metaData[action] && metaData[action].type === 'list' &&
            <Col span='10'>
              { metaData[action].children.map((item, index) => {
                if (item.type === 'select') {
                  return <FormItem key={index} label={item.name} {...formItemLayout}>
                    {getFieldDecorator(`${data.category}-${metaData[action].action}-${item.ci}-${data.uuid}`, {
                      initialValue: data.initialValue.metadata[item.ci],
                      rules: [{
                        required: item.required,
                        message: `${item.name}不能为空`
                      }]
                    })(
                      <Select style={{ width: '100%' }} placeholder={item.desc} disabled={disabled} >
                        {item.options.map((option, index) =>
                          <Option key={index} value={option.id}>{option.name}</Option>
                        )
                        }
                      </Select>
                    )}
                  </FormItem>;
                } else if (item.type === 'input') {
                  return <FormItem label={item.name} {...formItemLayout}>
                    {getFieldDecorator(`${data.category}-${metaData[action].action}-${item.ci}-${data.uuid}`, {
                      initialValue: data.initialValue.metadata[item.ci],
                      rules: [{
                        required: item.required,
                        message: `${item.name}不能为空`
                      }]
                    })(
                      <Input type={item.inputType} placeholder={item.desc} disabled={disabled} />
                    )}
                  </FormItem>;
                }

              })
              }
            </Col>
          }

          {
            (data.category === 'raid' || data.category === 'oob') && metaData[action] && metaData[action].type === 'input' &&
            <Col span='10' >
              <FormItem label={metaData[action].name} {...formItemLayout}>
                {getFieldDecorator(`${data.category}-${metaData[action].action}-${metaData[action].ci}-${data.uuid}`, {
                  initialValue: data.initialValue.metadata[metaData[action].ci],
                  rules: [{
                    required: metaData[action].required,
                    message: `${metaData[action].name}不能为空`
                  }]
                })(
                  <Input placeholder={metaData[action].desc} disabled={disabled} />
                )}
              </FormItem>
            </Col>
          }
          {
            (data.category === 'raid' || data.category === 'oob') && metaData[action] && metaData[action].type === 'radio' &&
            <Col span='10' >
              <FormItem label={metaData[action].name} {...formItemLayout}>
                {getFieldDecorator(`${data.category}-${metaData[action].action}-${metaData[action].ci}-${data.uuid}`,
                  { valuePropName: 'checked',
                    initialValue: data.initialValue.metadata[metaData[action].ci] === 'ON'
                  })(
                    <Switch checkedChildren='ON' unCheckedChildren='OFF' disabled={disabled} />
                )}
              </FormItem>
            </Col>
          }

          <Col span='3' className='colFixed'>
            {!disabled && <Button icon='copy' type='primary' size='small' style={{ marginRight: 8 }} onClick={() => this.props.copyConfig(data)} />}
            {!disabled && <Button icon='delete' type='danger' size='small' onClick={() => this.props.deleteConfig(index)} />}
          </Col>
        </div>
      </div>
    );
  }
}

export default Template;
