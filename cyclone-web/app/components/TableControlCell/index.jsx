import React from 'react';
import { Icon, Button, Tooltip, Menu, Dropdown } from 'antd';
import { Link } from 'react-router';

// eg
// const commands = [
//   {
//     name: "编辑",
//     command: "edit",
//     disabled: hasPermisson
//   },
//   {
//     name: "重启",
//     command: "restart",
//     disabled: hasPermisson
//   }
// ];

// return (
//   <TableControlCell
//     commands={commands}
//     record={record}
//     execCommand={command => {
//       this.execAction(command, [record]);
//     }}
//   />
// );

export default function Commands(props) {
  const commands = props.commands;
  const record = props.record;
  const limit = 3;
  const $commands = commands.slice(0, limit).map(renderCommand);
  let $more = null;
  if (commands.length > limit) {
    $commands.pop();
    $more = renderMore(commands.slice(limit - 1, commands.length));
  }
  $commands.push($more);
  return <div className='table-controls'>{$commands}</div>;

  /**
   * [renderMore description]
   * @param  {[type]} keys [description]
   * @return {[type]}      [description]
   */
  function renderMore(commands) {
    const $menu = (
      <Menu>
        {commands.map(command => {
          return (
            <Menu.Item key={`command-menu-${command.name}`}>
              {renderCommand(command)}
            </Menu.Item>
          );
        })}
      </Menu>
    );

    return (
      <Dropdown overlay={$menu} key='drodown'>
        <a className='ant-dropdown-link' href='javascript: void(0);'>
          更多 <Icon type='down' />
        </a>
      </Dropdown>
    );
  }

  /**
   * [renderCommand description]
   * @param  {[type]} key [description]
   * @return {[type]}     [description]
   */
  function renderCommand(command) {
    // disabled
    if (command.disabled) {
      return Disabled(command);
    }
    if (command.shouldDisabled && command.shouldDisabled(record)) {
      return Disabled(command);
    }
    // custom render
    if (command.render) {
      return command.render(props);
    }
    // default render
    return Default(command);
  }

  /**
   * [Default description]
   * @param {[type]} command [description]
   */
  function Default(command) {
    return (
      <a
        key={command.name}
        href='javascript:void(0)'
        style={{ color: command.type === 'danger' && '#ff3700' }}
        onClick={ev => {
          props.execCommand(command.command || command.name, record);
          ev.preventDefault();
        }}
      >
        {command.name}
      </a>
    );
  }

  /**
   * [Disabled description]
   * @param {[type]} command [description]
   * @param {[type]} message [description]
   */
  function Disabled(command) {
    return (
      <Tooltip key={command.name} title={command.message || command.name}>
        <span className='disabled' style={{ cursor: 'not-allowed' }}>{command.name}</span>
      </Tooltip>
    );
  }
}
