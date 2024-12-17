import React from 'react';
import { Row, Col } from 'antd';

export default class Layout extends React.Component {
  render() {
    return (
      <div className='approot-inner'>
        <div className='approot-inner-content'>
          <div className='page-body' style={{ paddingTop: 0 }}>
            <div className='page-tabs-unBordered'>
              <div className='page-unBorderedTitle'>{this.props.title}</div>
              <Row style={{ height: '100%' }}>
                <Col style={{ height: '100%' }}>
                  {this.props.children}
                </Col>
              </Row>
            </div>
          </div>
        </div>
      </div>
    );
  }
}
