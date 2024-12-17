import React from 'react';
import { post, get, getWithArgs, put } from 'common/xFetch2';
import {
  Form,
  Button,
  notification,
  Input,
  Radio,
  Select
} from 'antd';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout } from 'common/enums';
const Option = Select.Option;

class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.pxe_editor = null;
    this.sys_editor = null;
    this.state = {
      initialValue: {}
    };
  }

  componentDidMount() {
    this.initEditorValue();
  }

  componentWillUnmount() {
    this.pxe_editor = null;
    this.sys_editor = null;
  }

  initEditorValue = () => {
    const options = {
      lineNumbers: true,
      mode: 'shell',
      lineWrapping: false,
      readOnly: false,
      autoMatchParents: true,
      wordWrap: 'break-word',
      textWrapping: true,
      styleActiveLine: true,
      addModeClass: true,
      showCursorWhenSelecting: true
    };
    if ($('#pxeScript')[0]) {
      let editor = CodeMirror.fromTextArea($('#pxeScript')[0], { ...options, autofocus: true });
      if (this.props.type !== 'addSystem') {
        editor.setValue(this.props.initialValue.pxe);
      }
      editor.setSize('auto', 250);
      this.pxe_editor = editor;
    }
    if ($('#sysScript')[0]) {
      let editor = CodeMirror.fromTextArea($('#sysScript')[0], options);
      if (this.props.type !== 'addSystem') {
        editor.setValue(this.props.initialValue.content);
      }
      editor.setSize('auto', 250);
      this.sys_editor = editor;
    }
  };


  handleSubmit = ev => {
    ev && ev.preventDefault();
    ev && ev.stopPropagation();
    const { setFieldsValue } = this.props.form;
    const pxe_value = this.pxe_editor.getValue();
    const sys_value = this.sys_editor.getValue();
    setFieldsValue({ pxe: pxe_value, content: sys_value });
    this.props.form.validateFields({ force: true }, (error, values) => {
      if (error) {
        notification.warning({
          message: '还有未填写完整的选项'
        });
        return;
      }
      this.props.onSubmit(values);
    });
  };

  render() {
    const { osFamily, showSubmit, initialValue } = this.props;
    const { getFieldDecorator } = this.props.form;

    return (
      <div>
        <Form onSubmit={this.handleSubmit}>
          <FormItem {...formItemLayout} label='操作系统' hasFeedback={true}>
            {getFieldDecorator('family', {
              initialValue: initialValue.family,
              rules: [
                {
                  required: true,
                  message: '请选择操作系统'
                }
              ]
            })(<Select disabled={!showSubmit}>
              {
                !osFamily.loading &&
                osFamily.data.map((os, index) => <Option value={os.name} key={os.name}>{os.name}</Option>)
              }
            </Select>)}
          </FormItem>

          <FormItem {...formItemLayout} label='名称' hasFeedback={true}>
            {getFieldDecorator('name', {
              initialValue: initialValue.name,
              rules: [
                {
                  required: true,
                  message: '请输入名称'
                }
              ]
            })(
              <Input disabled={!showSubmit} />
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='启动模式'>
            {getFieldDecorator('boot_mode', {
              initialValue: initialValue.boot_mode,
              rules: [
                {
                  required: true,
                  message: '请选择启动模式'
                }
              ]
            })(
              <Radio.Group disabled={!showSubmit}>
                <Radio value='legacy_bios'>BIOS</Radio>
                <Radio value='uefi'>UEFI</Radio>
              </Radio.Group>
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='架构'>
            {getFieldDecorator('arch', {
              initialValue: initialValue.arch,
              rules: [
                {
                  required: true,
                  message: '请选择架构'
                }
              ]
            })(
              <Radio.Group disabled={!showSubmit}>
                <Radio value='x86_64'>x86_64</Radio>
                <Radio value='aarch64'>aarch64</Radio>
              </Radio.Group>
            )}
          </FormItem>           
          <FormItem {...formItemLayout} label='生命周期'>
            {getFieldDecorator('os_lifecycle', {
              initialValue: initialValue.os_lifecycle,
              rules: [
                {
                  required: true,
                  message: '请选择生命周期'
                }
              ]
            })(
              <Radio.Group disabled={!showSubmit}>
                <Radio value='testing'>Testing</Radio>
                <Radio value='active_default'>Active(Default)</Radio>
                <Radio value='active'>Active</Radio>
                <Radio value='containment'>Containment</Radio>
                <Radio value='end_of_life'>EOL</Radio>
              </Radio.Group>
            )}
          </FormItem> 
          <FormItem {...formItemLayout} label='用户名'>
            {getFieldDecorator('username', {
              initialValue: initialValue.username
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='密码'>
            {getFieldDecorator('password', {
              initialValue: initialValue.password
            })(<Input disabled={!showSubmit} type='password' />)}
          </FormItem>

          <FormItem {...formItemLayout} label='PXE模板内容'>
            {getFieldDecorator('pxe', {
              initialValue: initialValue.pxe,
              rules: [
                {
                  required: true,
                  message: '请填写PXE模板内容'
                }
              ]
            })(
              <div style={{ border: '1px solid #d9d9d9', marginBottom: '1px' }}>
                <textarea
                  disabled={!showSubmit}
                  id='pxeScript'
                  style={{ width: '100%', marginTop: '1px', marginBottom: '1px' }}
                />
              </div>
            )}
          </FormItem>

          <FormItem {...formItemLayout} label='系统模板内容'>
            {getFieldDecorator('content', {
              initialValue: initialValue.content,
              rules: [
                {
                  required: true,
                  message: '请填写系统模板内容'
                }
              ]
            })(
              <div style={{ border: '1px solid #d9d9d9', marginBottom: '1px' }}>
                <textarea
                  disabled={!showSubmit}
                  id='sysScript'
                  style={{ width: '100%', marginTop: '1px', marginBottom: '1px' }}
                />
              </div>
            )}
          </FormItem>
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='创建时间'>
              {getFieldDecorator('created_at', {
                initialValue: initialValue.created_at
              })(<Input disabled={!showSubmit} />)}
            </FormItem>
          }
          {
            !showSubmit &&
            <FormItem {...formItemLayout} label='修改时间'>
              {getFieldDecorator('updated_at', {
                initialValue: initialValue.updated_at
              })(<Input disabled={!showSubmit} />)}
            </FormItem>
          }

          <FormItem {...tailFormItemLayout}>
            <div className='pull-right'>
              <Button onClick={() => this.props.onCancel()}>取消</Button>
              {
                showSubmit &&
                <Button
                  style={{ marginLeft: 8 }}
                  type='primary'
                  htmlType='submit'
                >
                  提交
                </Button>
              }
            </div>
          </FormItem>
        </Form>
      </div>
    );
  }
}

export default Form.create()(MyForm);
