import React from 'react';
import { Button, Col, Row, Form, Input, Select, Icon } from 'antd';
const FormItem = Form.Item;
import { formItemLayout, tailFormItemLayout, tableFormItemLayout } from 'common/enums';
const Option = Select.Option;

export default class Disks extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      disks: [],
      loading: true
    };
  }
  componentDidMount() {
    const { disks } = this.props;
    if (disks && disks.length > 0) {
      this.setState({
        disks: disks,
        loading: false
      });
    }
  }
  addDisks = () => {
    this.setState({
      disks: this.props.initDisks,
      loading: false
    });
    const { setFieldsValue } = this.props.form;
    setFieldsValue({ disks: this.props.initDisks });
  };
  deleteDisk = (index) => {
    const { disks } = this.state;
    disks.splice(index, 1);
    this.setState({ disks });
  };
  copyDisk = (data) => {
    const { disks } = this.state;
    const copyData = JSON.parse(JSON.stringify(data));
    disks.push(copyData);
    this.setState({ disks });
  };
  copyDiskPart = (data, index) => {
    const { disks } = this.state;
    const copyData = JSON.parse(JSON.stringify(data));
    disks[index].partitions.push(copyData);
    this.setState({ disks });
  };
  deleteDiskPart = (part_index, index) => {
    const { disks } = this.state;
    disks[index].partitions.splice(part_index, 1);
    this.setState({ disks });
  };
  render() {
    const { disabled, form } = this.props;
    const { disks, loading } = this.state;
    const { getFieldDecorator } = form;
    return (
      <div>
        {
          !disks || (disks || []).length === 0 &&
          <Button icon='plus' onClick={() => this.addDisks()}>添加</Button>
        }
        {
            !loading && (disks || []).map((item, index) => {
              return (
                <div className='disk-panel'>
                  <div className='disk'>
                    <Row className='disk-header'>
                      <Col span='20' >
                        <FormItem {...formItemLayout} label='磁盘名称'>
                          {getFieldDecorator(`diskParams-diskName-${index}`, {
                            initialValue: item.name,
                            rules: [{
                              required: true,
                              message: ' '
                            }]
                          })(<Input placeholder='请填写磁盘名称' disabled={disabled} />)}
                        </FormItem>
                      </Col>
                      <Col span='4' className='col-right'>
                        {!disabled &&
                        <Icon type='copy' theme='outlined' style={{ marginRight: 8, color: '#108ee9' }} onClick={() => this.copyDisk(item)} />
                      }
                        {!disabled && index !== 0 &&
                        <Icon type='close' theme='outlined' style={{ color: '#ff3700' }} onClick={() => this.deleteDisk(index)} />}
                      </Col>
                    </Row>
                    <Row className='disk-table-header'>
                      <Col span='5'>名称</Col>
                      <Col span='5'>大小</Col>
                      <Col span='5'>文件系统类型</Col>
                      <Col span='5'>挂载点</Col>
                      <Col span='4' />
                    </Row>
                    {
                    (item.partitions || []).map((part, part_index) => {
                      return (
                        <Row className='disk-table-body'>
                          <Col span='5'>
                            <FormItem {...tableFormItemLayout}>
                              {getFieldDecorator(`diskParams-name-${index}-${part_index}`, {
                                initialValue: part.name,
                                rules: [{
                                  required: true,
                                  message: ' '
                                }]
                              })(<Input placeholder='请填写名称' disabled={disabled} />)}
                            </FormItem>
                          </Col>
                          <Col span='5'>
                            <FormItem {...tableFormItemLayout}>
                              {getFieldDecorator(`diskParams-size-${index}-${part_index}`, {
                                initialValue: part.size,
                                rules: [{
                                  required: true,
                                  message: ' '
                                }]
                              })(<Input placeholder='请填写大小' disabled={disabled} />)}
                            </FormItem>
                          </Col>
                          <Col span='5'>
                            <FormItem {...tableFormItemLayout}>
                              {getFieldDecorator(`diskParams-fstype-${index}-${part_index}`, {
                                initialValue: part.fstype,
                                rules: [{
                                  required: true,
                                  message: ' '
                                }]
                              })(<Select placeholder='请选择文件系统类型' disabled={disabled}>
                                <Option value='ext3'>ext3</Option>
                                <Option value='ext4'>ext4</Option>
                                <Option value='xfs'>xfs</Option>
                                <Option value='swap'>swap</Option>
                                <Option value='ntfs'>ntfs</Option>
                                <Option value='vfat'>vfat</Option>
                              </Select>)}
                            </FormItem>
                          </Col>
                          <Col span='5'>
                            <FormItem {...tableFormItemLayout}>
                              {getFieldDecorator(`diskParams-mountpoint-${index}-${part_index}`, {
                                initialValue: part.mountpoint,
                                rules: [{
                                  required: true,
                                  message: ' '
                                }]
                              })(<Input placeholder='请填写挂载点' disabled={disabled} />)}
                            </FormItem>
                          </Col>
                          <Col span='4' className='col-right'>
                            {!disabled && <Button size='small' icon='copy' style={{ marginRight: 8, color: '#108ee9' }} onClick={() => this.copyDiskPart(part, index)} />}
                            {!disabled && <Button type='danger' disabled={part_index === 0} size='small' icon='delete' onClick={() => this.deleteDiskPart(part_index, index)} />}
                          </Col>
                        </Row>
                      );
                    })
                  }
                  </div>
                </div>
              );
            })
          }
      </div>
    );
  }
}
