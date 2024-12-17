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
import Disk from './disks';
const Option = Select.Option;


class MyForm extends React.Component {
  constructor(props) {
    super(props);
    this.pre_editor = null;
    this.post_editor = null;
    this.state = {
      initialValue: {},
      disks: []
    };
  }
  componentWillUnmount() {
    this.pre_editor = null;
    this.post_editor = null;
  }

  componentDidMount() {
    this.initEditorValue();
  }


  initEditorValue = () => {
    const options = {
      lineNumbers: true,
      mode: 'shell',
      readOnly: false,
      autoMatchParens: true,
      wordWrap: 'break-word',
      textWrapping: true,
      styleActiveLine: true,
      scrollbarStyle: 'native',
      addModeClass: true,
      showCursorWhenSelecting: true
    };
    if ($('#preScript')[0]) {
      let editor = CodeMirror.fromTextArea($('#preScript')[0], options);
      if (this.props.type !== 'addMirror') {
        editor.setValue(this.props.initialValue.pre_script);
      }
      editor.setSize('auto', 250);
      this.pre_editor = editor;
    }
    if ($('#postScript')[0]) {
      let editor = CodeMirror.fromTextArea($('#postScript')[0], options);
      editor.setSize('auto', 250);
      if (this.props.type !== 'addMirror') {
        editor.setValue(this.props.initialValue.post_script);
      }
      this.post_editor = editor;
    }
  };

  handleSubmit = ev => {
    ev && ev.preventDefault();
    ev && ev.stopPropagation();

    const { setFieldsValue } = this.props.form;
    const pre_value = this.pre_editor.getValue();
    const post_value = this.post_editor.getValue();
    setFieldsValue({ pre_script: pre_value, post_script: post_value });

    this.props.form.validateFields({ force: true }, (error, values) => {

      const diskKeys = Object.keys(values);
      const fieldParamsKeys = diskKeys.filter(
        key => key.indexOf('diskParams-diskName') >= 0
      );
      const disks = fieldParamsKeys.map(key => {
        const id = key.replace('diskParams-diskName-', '');
        const sizeParamsKeys = diskKeys.filter(
          key => key.indexOf(`diskParams-size-${id}`) >= 0
        );

        const partitions = sizeParamsKeys.map(sizeKey => {
          const partId = sizeKey.replace(`diskParams-size-${id}-`, '');
          return {
            size: values[`diskParams-size-${id}-${partId}`],
            fstype: values[`diskParams-fstype-${id}-${partId}`],
            mountpoint: values[`diskParams-mountpoint-${id}-${partId}`],
            name: values[`diskParams-name-${id}-${partId}`]
          };
        });
        return {
          name: values[`diskParams-diskName-${id}`],
          partitions: partitions
        };
      });

      const data = {
        family: values.family,
        boot_mode: values.boot_mode,
        arch: values.arch,
        os_lifecycle: values.os_lifecycle,
        name: values.name,
        url: values.url,
        username: values.username,
        password: values.password,
        disks: disks,
        pre_script: values.pre_script,
        post_script: values.post_script
      };

      if (error) {
        notification.warning({
          message: '还有未填写完整的选项'
        });
        return;
      }
      this.props.onSubmit(data);
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
          <FormItem {...formItemLayout} label='镜像下载地址' hasFeedback={true}>
            {getFieldDecorator('url', {
              initialValue: initialValue.url,
              rules: [
                {
                  required: true,
                  message: '请输入镜像下载地址'
                }
              ]
            })(
              <Input disabled={!showSubmit} />
            )}
          </FormItem>
          <FormItem {...formItemLayout} label='用户名' hasFeedback={true}>
            {getFieldDecorator('username', {
              initialValue: initialValue.username
            })(<Input disabled={!showSubmit} />)}
          </FormItem>
          <FormItem {...formItemLayout} label='密码' hasFeedback={true}>
            {getFieldDecorator('password', {
              initialValue: initialValue.password
            })(<Input disabled={!showSubmit} type='password' />)}
          </FormItem>
          <FormItem {...formItemLayout} label='分区方案'>
            {getFieldDecorator('disks', {
              initialValue: initialValue.disks,
              rules: [
                {
                  required: true,
                  message: '分区方案不能为空'
                }
              ]
            })(
              <Disk initDisks={this.props.initDisks} disabled={!showSubmit} disks={initialValue.disks} form={this.props.form} />
            )}
          </FormItem>

          <FormItem {...formItemLayout} label='前置脚本'>
            {getFieldDecorator('pre_script', {
              initialValue: initialValue.pre_script
            })(
              <div style={{ border: '1px solid #d9d9d9', marginBottom: '1px' }}>
                <textarea
                  disabled={!showSubmit}
                  id={'preScript'}
                  style={{ width: '100%', marginTop: '1px', marginBottom: '1px' }}
                />
              </div>
            )}
          </FormItem>

          <FormItem {...formItemLayout} label='后置脚本'>
            {getFieldDecorator('post_script', {
              initialValue: initialValue.post_script
            })(
              <div style={{ border: '1px solid #d9d9d9', marginBottom: '1px' }}>
                <textarea
                  disabled={!showSubmit}
                  id={'postScript'}
                  style={{ width: '100%', marginTop: '1px', marginBottom: '1px' }}
                />
              </div>
            )}
          </FormItem>

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
