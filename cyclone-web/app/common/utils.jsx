import React from 'react';
import { Tooltip, Icon, Row, Col, Breadcrumb, Checkbox, Popover, Button } from 'antd';
import { Link } from 'react-router';
const CheckboxGroup = Checkbox.Group;
import { arrayFind } from './util';

/**
 *
 * @param label
 * @param title
 * @returns {XML}
 */
export const tpsLabel = (label, title) => {
  if (title) {
    return (
      <span>
        {label}
        <Tooltip title={title}>
          <Icon type='question-circle-o' />
        </Tooltip>
      </span>
    );
  }
  return <span>{label}</span>;
};

/**
 *
 * @param fields
 * @param span
 */
export const renderFormDetail = (fields, span = 8) => {
  return fields.map(node => {
    return (
      <Col span={span}>
        <span className='panel-label'>
          {node.label}：
        </span>
        <span className='panel-value'>
          {node.value}
        </span>
      </Col>
    );
  });
};

/**
 *
 * @param title
 * @param link
 * @returns {XML}
 */
export const getBreadcrumb = (title, link) => {
  return (
    <Breadcrumb>
      <Breadcrumb.Item>
        {link ?
          <Link to={link}><Icon type='left' />返回</Link> :
          <a onClick={() => history.back()}><Icon type='left' />返回</a>
        }
      </Breadcrumb.Item>
      <Breadcrumb.Item>{title}</Breadcrumb.Item>
    </Breadcrumb>
  );
};

/**
 *
 * @param title
 * @returns {XML}
 */
export const geTabsTitle = (title) => {
  return (
    <Breadcrumb>
      <Breadcrumb.Item><h3>{title}</h3></Breadcrumb.Item>
    </Breadcrumb>
  );
};

/**
 * 显示字段
 * @param self
 * @returns {XML}
 */
export const renderDisplayMore = (self, plainOptions = []) => {
  const onChange = checkedList => {
    self.setState({
      checkedList,
      indeterminate:
      !!checkedList.length && checkedList.length < plainOptions.length,
      checkAll: checkedList.length === plainOptions.length
    });
  };
  const onCheckAllChange = e => {
    self.setState({
      checkedList: e.target.checked ? plainOptions : [],
      indeterminate: false,
      checkAll: e.target.checked
    });
  };
  return (
    <Popover
      placement='bottom'
      content={
        <div>
          <div style={{ borderBottom: '1px solid #E9E9E9' }}>
            <Checkbox
              indeterminate={self.state.indeterminate}
              onChange={onCheckAllChange}
              checked={self.state.checkAll}
            >
              全选
            </Checkbox>
          </div>
          <CheckboxGroup value={self.state.checkedList} onChange={onChange}>
            {plainOptions.map(data => {
              return (
                <Row key={data.name} >
                  <Col span={24}>
                    <Checkbox value={data}>{data.name}</Checkbox>
                  </Col>
                </Row>
              );
            })}
          </CheckboxGroup>
        </div>
      }
      trigger='click'
    >
      <Button>
        显示字段<Icon type='down' />
      </Button>
    </Popover>
  );
};

export const getPermissonBtn = (permissions = [], btn) => {
  return arrayFind(permissions, btn);
};
