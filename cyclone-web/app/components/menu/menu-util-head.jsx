/**
 * Created by souakiragen on 2017/6/5.
 */
import React, { Component } from 'react';
import { Menu, Icon, Dropdown } from 'antd';
import { getPermissonBtn } from 'common/utils';

export default class PermissionMenus extends Component {

  render() {
    const menu = (
      <Menu>
        <Menu.Item>
          <a
            onClick={() => {
              this.props.showAbout();
            }}
          >
            关于
          </a>
        </Menu.Item>
        <Menu.Item>
          <a
            onClick={() => {
              this.props.changePwd();
            }}
          >
            修改密码
          </a>
        </Menu.Item>
        <Menu.Item>
          <a
            style={{ borderTop: '1px solid #e2e7ec' }}
            onClick={() => {
              this.props.logout();
            }}
          >
        退出登录
        </a>
        </Menu.Item>
      </Menu>
    );
    const userData = this.props.userData;
    return (
      <div>
        <div className='app-header-btn' onClick={this.props.handlerCollapse}>
          <div className='btn-circle'>
            <i className={'conf icon-Group1'} />
          </div>
        </div>
        <div className='app-header-left'>
          <div className='app-header-logo'>
            <img className='logo' src={this.props.logoUrl} />
          </div>
        </div>
        <div className='app-header-login'>
          <span
            className='app-theme'
            style={{ marginRight: '10px' }}
            onClick={this.props.changeTheme}
          >
            <i
              className={'conf icon-painter-palette'}
              style={{ marginRight: 8 }}
            />
          </span>
          <Dropdown overlay={menu} style={{ float: 'right' }}>
            <a className='ant-dropdown-link'>
              {userData.name} <Icon type='down' />
              <img src='assets/icon/icon_avatar.png' alt='' width='32' style={{ marginLeft: 8 }}/>
            </a>
          </Dropdown>
        </div>
      </div>
    );
  }
}
