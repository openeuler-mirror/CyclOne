/**
 * App Component
 */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Icon, Spin } from 'antd';
import pureRender from 'common/pureRender';
import PermissionMenusHead from 'components/menu/menu-util-head';
import classnames from 'classnames';
import PermissionMenusMain from 'components/menu/menu-util-main';
import PermissionMenusMainCollapsed from 'components/menu/menu-util-main-collapsed';
import { post } from 'common/xFetch2';
import Popup from 'components/popup';
import changePassword from './changePassword';

/**
 * [createClass description]
 * @param  {[type]} { render( [description]
 * @return {[type]}            [description]
 */

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      ifDefaultTheme: true,
      collapsed: false,
      openKeys: [],
      lastMenu: '',
      logoUrl: '/assets/logo/logo-cyclone.png',
      menuKeyMap: {
        'database/idc': 'database/idc',
        'database/room': 'database/room',
        'database/network': 'database/network',
        'database/cabinet': 'database/cabinet',
        'database/usite': 'database/usite',
        'database/store-room': 'database/store-room',
        'network/cidr': 'network/cidr',
        'network/ips': 'network/ips',
        'network/device': 'network/device',
        'device/list': 'device/list',
        'device/detail': 'device/list',
        'device/oob': 'device/oob',
        'device/special': 'device/special',
        'device/pre_deploy': 'device/pre_deploy',
        'device/setting': 'device/setting',
        'device/setting_rule': 'device/setting_rule',
        'device/entry': 'device/pre_deploy',
        'device/inspection/list': 'device/inspection/list',
        'device/inspection/detail': 'device/inspection/list',
        'template/system': 'template/system',
        'template/hardware/list': 'template/hardware/list',
        'template/hardware/create': 'template/hardware/list',
        'template/hardware/detail': 'template/hardware/list',
        'template/hardware/edit': 'template/hardware/list',
        'audit/log': 'audit/log',
        'audit/api': 'audit/api',
        'approval': 'approval',
        'order/list': 'order/list',
        'order/deviceCategory': 'order/deviceCategory',
        'task/list': 'task/list'
      }
    };
  }

  componentDidMount() {
    this.checkLogin();
    this.getUserData();
    this.props.dispatch({
      type: 'global/idc/getSearchList'
    });
    this.props.dispatch({
      type: 'global/room/getSearchList'
    });
    this.props.dispatch({
      type: 'global/cabinet/getSearchList'
    });
    this.props.dispatch({
      type: 'global/networkArea/getSearchList'
    });
  }

  changePwd = () => {
    changePassword({
      reload: this.logout
    });
  };

  showAbout = () => {
    Popup.open({
      title: '版本信息',
      width: 416,
      onCancel: () => {
        Popup.close();
      },
      content: (
        <div className='ant-confirm-body-wrapper'>
          <div className='ant-confirm-body' style={{ paddingBottom: 26 }}>
            <img src='assets/icon/about_icon.png' alt='' width='40' style={{ marginRight: 16, float: 'left' }}/>
            <span className='ant-confirm-title'>当前版本号 V1.5.5</span>
            {/*<div className='ant-confirm-content'>*/}
            {/*<p>请尽快升级</p>*/}
            {/*</div>*/}
          </div>
        </div>
      )
    });
  };

  logout = () => {
    this.props.dispatch({
      type: 'global/logout'
    });
  };

  checkLogin = () => {
    this.props.dispatch({
      type: 'global/checkLogin'
    });
  };

  getUserData = () => {
    this.props.dispatch({
      type: 'global/getUserData'
    });
  };

  // 展开或收起菜单
  handlerCollapse = () => {
    this.setState({
      collapsed: !this.state.collapsed
    });
  };

  changeTheme = () => {
    this.setState({
      ifDefaultTheme: !this.state.ifDefaultTheme
    });
  };

  handlerOpenChange = openKeys => {
    if (openKeys.length > 0) {
      this.setState({
        openKeys: new Array(openKeys[openKeys.length - 1])
      });
    } else {
      this.setState({
        openKeys: []
      });
    }
  };

  getOpenKeys = () => {
    return this.state.openKeys;
  };
  getIconStation(key) {
    const openKey = this.getOpenKeys();
    const collapsed = this.state.collapsed;
    if (openKey[0] === key && collapsed === true) {
      return true;
    } else {
      return false;
    }
  }
  getCurrentMenu = () => {
    let path = window.location.hash.match(/#\/(.+)/),
      ret = '/';
    path = path ? path[1] : '';

    if (path == '') {
      return ret;
    }

    if (path == 'login') {
      return path;
    }
    const menuKeyMap = this.state.menuKeyMap;
    Object.keys(menuKeyMap).forEach(k => {
      const v = k;
      if (path.indexOf(v) == 0) {
        ret = menuKeyMap[k];
      }
    });
    return ret;
  };
  render() {
    const current = this.getCurrentMenu();
    const { userData, login } = this.props.data;
    const permissions = userData ? userData.permissions : [];
    const menuItems = [
      {
        key: '/',
        permissionKey: 'menu_home',
        link: '/',
        icon: <Icon type='home' />,
        description: '概览'
      },
      {
        key: 'database',
        permissionKey: 'menu_idc_management',
        icon:
          this.getIconStation('database') === true ? (
            <Icon type='down' />
          ) : (
            <Icon type='database' />
          ),
        description: '数据中心',
        children: [
          {
            key: 'database/idc',
            permissionKey: 'menu_idc',
            link: '/database/idc',
            description: '数据中心管理',
            icon: (
              <Icon type='global' theme='outlined' />
            )
          },
          {
            key: 'database/room',
            permissionKey: 'menu_server_room',
            link: '/database/room',
            description: '机房信息管理',
            icon: (
              <Icon type='bank' theme='outlined' />)
          },
          {
            key: 'database/network',
            permissionKey: 'menu_network_area',
            link: 'database/network',
            description: '网络区域管理',
            icon: (
              <Icon type='wifi' theme='outlined' />)
          },
          {
            key: 'database/cabinet',
            permissionKey: 'menu_server_cabinet',
            link: 'database/cabinet',
            description: '机架信息管理',
            icon: (
              <Icon type='cluster' theme='outlined' />)
          },
          {
            key: 'database/usite',
            permissionKey: 'menu_server_usite',
            link: 'database/usite',
            description: '机位信息管理',
            icon: (
              <Icon type='gold' theme='outlined' />)
          },
          {
            key: 'database/store-room',
            permissionKey: 'menu_store_room',
            link: 'database/store-room',
            description: '库房信息管理',
            icon: (
              <Icon type='shop' theme='outlined' />)
          }
        ]
      },
      {
        key: 'network',
        permissionKey: 'menu_network_management',
        icon:
          this.getIconStation('network') === true ? (
            <Icon type='down' />
          ) : (
            <Icon type='share-alt' />
          ),
        description: '网络管理',
        children: [
          {
            key: 'network/device',
            permissionKey: 'menu_network_device',
            link: '/network/device',
            description: '网络设备管理',
            icon: (
              <Icon type='laptop' theme='outlined' />
            )
          },
          {
            key: 'network/cidr',
            permissionKey: 'menu_ip_network',
            link: '/network/cidr',
            description: 'IP网段管理',
            icon: (
              <Icon type='build' theme='outlined' />
            )
          },
          {
            key: 'network/ips',
            permissionKey: 'menu_ip',
            link: '/network/ips',
            description: 'IP分配管理',
            icon: (
              <Icon type='gateway' theme='outlined' />
            )
          }
        ]
      },
      {
        key: 'device',
        permissionKey: 'menu_physical_machine_management',
        icon:
          this.getIconStation('device') === true ? (
            <Icon type='down' />
          ) : (
            <Icon type='desktop' />
          ),
        description: '物理机管理',
        children: [
          {
            key: 'device/list',
            permissionKey: 'menu_physical_machine',
            link: '/device/list',
            description: '物理机列表',
            icon: (
              <Icon type='project' theme='outlined' />
            )
          },
          {
            key: 'device/oob',
            permissionKey: 'menu_oob_info',
            link: '/device/oob',
            description: '带外管理',
            icon: (
              <Icon type='wifi' />
            )
          },
          {
            key: 'device/special',
            permissionKey: 'menu_special_device',
            link: '/device/special',
            description: '特殊设备',
            icon: (
              <Icon type='block' />
            )
          },
          {
            key: 'device/pre_deploy',
            permissionKey: 'menu_predeploy_physical_machine',
            link: '/device/pre_deploy',
            description: '待部署物理机',
            icon: (
              <Icon type='exception' theme='outlined'/>
            )
          },
          {
            key: 'device/setting',
            permissionKey: 'menu_device_setting',
            link: '/device/setting',
            description: '部署列表',
            icon: (
              <Icon type='cluster' theme='outlined' />
            )
          },
          {
            key: 'device/setting_rule',
            permissionKey: 'menu_device_setting_rule',
            link: '/device/setting_rule',
            description: '部署参数规则',
            icon: (
              <Icon type='cluster' theme='outlined' />
            )
          },          
          {
            key: 'device/inspection/list',
            permissionKey: 'menu_inspection',
            link: '/device/inspection/list',
            description: '硬件巡检',
            icon: <Icon type='tool' />
          }
        ]
      },
      {
        key: 'order',
        permissionKey: 'menu_order_management',
        icon:
          this.getIconStation('order') === true ? (
            <Icon type='down' />
          ) : (
            <Icon type='profile' />
          ),
        description: '订单管理',
        children: [
          {
            key: 'order/list',
            permissionKey: 'menu_order',
            link: '/order/list',
            description: '订单列表',
            icon: (
              <Icon type='project' theme='outlined' />
            )
          },
          {
            key: 'order/deviceCategory',
            permissionKey: 'menu_device_category',
            link: '/order/deviceCategory',
            description: '设备类型',
            icon: (
              <Icon type='crown' />
            )
          }
        ]
      },
      {
        key: 'template',
        permissionKey: 'menu_template_management',
        icon: this.getIconStation('template') === true ? <Icon type='down' /> : <Icon type='medicine-box' theme='outlined' />,
        description: '配置管理',
        children: [
          {
            key: 'template/system',
            permissionKey: 'menu_system_template',
            link: '/template/system',
            description: '装机配置',
            icon: (
              <Icon type='code' theme='outlined' />
            )
          },
          {
            key: 'template/hardware/list',
            permissionKey: 'menu_hardware_template',
            link: '/template/hardware/list',
            description: '硬件配置',
            icon: <Icon type='setting' />
          }
        ]
      },
      {
        key: 'approval',
        permissionKey: 'menu_approval',
        link: 'approval',
        icon: <Icon type='safety-certificate' />,
        description: '审批管理'
      },
      { key: 'task/list',
        permissionKey: 'menu_task_management',
        link: 'task/list',
        icon: <Icon type='schedule' />,
        description: '任务管理'
      },
      {
        key: 'audit',
        permissionKey: 'menu_audit',
        icon: this.getIconStation('template') === true ? <Icon type='down' /> : <Icon type='audit' />,
        description: '操作审计',
        children: [
          {
            key: 'audit/log',
            permissionKey: 'menu_audit_log',
            link: '/audit/log',
            description: '操作记录',
            icon: (
              <Icon type='file-text' />
            )
          },
          {
            key: 'audit/api',
            permissionKey: 'menu_audit_api',
            link: '/audit/api',
            description: '接口调用记录',
            icon: (
              <Icon type='file-sync' />
            )
          }
        ]
      },
      {
        key: 'user',
        linkGo: `${userData.uam_portal_url}`,
        permissionKey: 'menu_user_management',
        iconGo: <Icon type='user' theme='outlined' />,
        description: '用户管理'
      }
    ];
    const style = {
      width: this.state.collapsed ? '50px' : '200px',
      minWidth: this.state.collapsed ? '50px' : '200px'
    };
    const openKey = this.getOpenKeys();
    const menuPreRender = (filterMenuItems = []) => {
      if (filterMenuItems && filterMenuItems.length > 0) {
        let path = window.location.hash.match(/#\/(.+)\??/);
        path = path ? path[1] : '';
        if (
          menuItems.some(item => item.link === '/' + path) &&
          !filterMenuItems.some(item => item.link === '/' + path)
        ) {
          window.location.href = '/#' + filterMenuItems[0].link;
        }
      }
    };
    return (
      <div
        className={classnames({
          approot: true,
          themeDefault: this.state.ifDefaultTheme,
          themeDark: !this.state.ifDefaultTheme
        })}
      >
        {
          login === false && (
            <div
              style={{
                textAlign: 'center',
                position: 'fixed',
                background: '#fff',
                top: 0,
                bottom: 0,
                left: 0,
                right: 0,
                zIndex: 9999
              }}
            >
              <Spin size='large' tip='还没有登录' style={{ fontSize: 20, height: '100%', transform: 'translateY(50%)' }} />
            </div>
          )
        }
        <div>
          <div className='app-header'>
            <PermissionMenusHead
              logoUrl={this.state.logoUrl}
              mode='inline'
              theme='dark'
              userData={this.props.data.userData}
              dispatch={this.props.dispatch}
              logout={this.logout}
              changePwd={this.changePwd}
              showAbout={this.showAbout}
              handlerCollapse={this.handlerCollapse}
              changeTheme={this.changeTheme}
            />
          </div>
          <div className='app-container'>
            <div
              className={classnames({
                'app-body-nav': true,
                'app-body-nav-collapsed': this.state.collapsed,
                'app-body-nav-uncollapsed': !this.state.collapsed
              })}
              style={style}
            >
              {!this.state.collapsed && (
              <PermissionMenusMain
                theme='dark'
                menuItems={menuItems}
                permissions={permissions}
                dispatch={this.props.dispatch}
                mode='inline'
                onOpenChange={this.handlerOpenChange}
                openKeys={openKey}
                selectedKeys={[current]}
                menuPreRender={menuPreRender}
              />
                )}
              {this.state.collapsed && (
              <PermissionMenusMainCollapsed
                theme='dark'
                userData={this.props.data.userData}
                menuItems={menuItems}
                permissions={permissions}
                dispatch={this.props.dispatch}
                mode='inline'
                selectedKeys={[current]}
                menuPreRender={menuPreRender}
                onOpenChange={this.handlerOpenChange}
                openKeys={openKey}
              />
                )}
              <div
                className='app-body-nav-btn'
                onClick={this.handlerCollapse}
              >
                {!this.state.collapsed && (
                <div className='app-body-nav-uncollapsed'>
                  <img src='/assets/logo/logo-max.png' alt='' height='30' />
                </div>
                  )}
                {this.state.collapsed && (
                <div className='app-body-nav-collapsed'>
                  <img width='20' src='/assets/logo/logo-min.png' alt='' />
                </div>
                  )}
              </div>
            </div>
            <div className='app-body-main app-body__project'>
              <div className='app-body-content'> {React.Children.map(this.props.children,
              (child) => React.cloneElement(child, { permissions }))}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
}

/**
 * transform store data to App props
 * @param  {[type]} state [description]
 * @return {[type]}       [description]
 */
function mapStateToProps(store) {
  return {
    data: store.get('global').toJS()
  };
}

/**
 * add dispatch to
 * @param  {[type]} dispatch [description]
 * @return {[type]}          [description]
 */
function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

/**
 * create redux container
 * @type {[type]}
 */

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(pureRender(App));
