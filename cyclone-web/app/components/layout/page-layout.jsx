import React from 'react';

export default class Layout extends React.Component {
  render() {
    return (
      <div className='approot-inner'>
        <div className='approot-inner-content'>
          <div className={'page-node-report page'}>
            {this.props.title && <div className='page-header'>{this.props.title}</div>}
            <div className='page-body'>
              {this.props.children}
            </div>
          </div>
        </div>
      </div>
    );
  }
}
