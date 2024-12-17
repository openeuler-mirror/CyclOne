/**
 * Created by souakiragen on 2017/6/5.
 */
import React, { Component } from 'react';
import { arrayFind } from '../../common/util';
import { IndexLink, hashHistory } from 'react-router';
import { Menu, Icon, notification, Dropdown, Tooltip } from 'antd';
const SubMenu = Menu.SubMenu;

export default class PermissionMenus extends Component {
  render() {
    const renderIcon = item => {
      if (item.transform) {
        return item.transform(item.icon);
      } else {
        return <Icon type={item.icon} />;
      }
    };


    const renderTitle = item => {
      if (item.titleFormat) {
        return item.titleFormat();
      } else {
        return (
          <span>
            {item.icon} <span>{item.description}</span>
          </span>
        );
      }
    };
    const renderMenuItem = item => {
      if (item.children) {
        return (
          <SubMenu key={item.key} title={renderTitle(item)}>
            {item.children.map(child => {
              return renderMenuItem(child);
            })}
          </SubMenu>
        );
      } else if (item.icon) {
        return (
          <Menu.Item key={item.key} disabled={item.disabled}>
            <IndexLink to={item.link}>
              {item.icon} <span>{item.description}</span>
            </IndexLink>
          </Menu.Item>
        );

      } else if (item.linkGo) {
        return (
          <Menu.Item key={item.key} disabled={item.disabled}>
            <a href={item.linkGo} target='blank'>
              {item.iconGo} <span>{item.description}</span>
            </a>
          </Menu.Item>
        );
      } else {
        return (
          <Menu.Item key={item.key} disabled={item.disabled}>
            <IndexLink to={item.link}>{item.description}</IndexLink>
          </Menu.Item>
        );
      }
    };

    const permissions = this.props.permissions;
    const filterMenuItems = (items = []) => {
      const results = items.filter(item => {
        return (
          item.withPermission || arrayFind(permissions, item.permissionKey)
        );
      });
      for (const result of results) {
        if (result.children) {
          result.children = filterMenuItems(result.children);
        }
      }
      return results;
    };
    const menuItems = filterMenuItems(this.props.menuItems);

    if (this.props.menuPreRender) {
      this.props.menuPreRender(menuItems);
    }

    const menus = menuItems.map(item => {
      return renderMenuItem(item);
    });

    return <Menu {...this.props}>{menus}</Menu>;
  }
}
