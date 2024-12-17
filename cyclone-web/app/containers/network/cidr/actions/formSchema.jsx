import React from 'react';
import { NET_STATUS, getSearchList } from 'common/enums';
import { IP_NETWORK_CATEGORY, IP_VERSION } from "../../../../common/enums";

export const formSchema = (disabled, options, type) => {
  let elements = [
    {
      id: 'cidr',
      name: 'cidr',
      label: '网段名称',
      type: 'TextInput',
      disabled: disabled || type === '_update',
      rules: [
        {
          required: true,
          message: '请填写网段名称'
        }
      ]
    },
    {
      id: 'server_room_id',
      name: 'server_room_id',
      label: '机房管理单元',
      type: 'Select',
      disabled: disabled || type === '_update',
      rules: [
        {
          required: true,
          message: '请选择机房管理单元'
        }
      ],
      options: !options.room.loading ? (options.room.data || []).map(it => {
        return {
          label: it.name, value: it.id, key: it.id
        };
      }) : []
    },
    {
      id: 'category',
      name: 'category',
      label: '网段类别',
      type: 'Select',
      disabled: disabled || type === '_update',
      rules: [
        {
          required: true,
          message: '请选择网段类别'
        }
      ],
      options: getSearchList(IP_NETWORK_CATEGORY)
    },
    {
      id: 'netmask',
      name: 'netmask',
      label: '掩码',
      type: 'TextInput',
      disabled: disabled || type === '_update',
      rules: [
        {
          required: true,
          message: '请输入掩码'
        }
      ]
    },
    {
      id: 'gateway',
      name: 'gateway',
      label: '网关',
      type: 'TextInput',
      disabled: disabled || type === '_update',
      rules: [
        {
          required: true,
          message: '请输入网关'
        }
      ]
    },
    {
      id: 'ip_pool',
      name: 'ip_pool',
      label: 'IP资源池',
      type: 'TextInput',
      disabled: disabled || type === '_update',
      placeholder: '用,分隔',
      rules: [
        {
          required: true,
          message: '请输入IP资源池'
        }
      ]
    },
    {
      id: 'pxe_pool',
      name: 'pxe_pool',
      label: 'PXE资源池',
      type: 'TextInput',
      placeholder: '用,分隔',
      disabled: disabled || type === '_update',
      rules: [
        {
          message: '请输入PXE资源池'
        }
      ]
    },
    {
      id: 'switchs',
      name: 'switchs',
      label: '覆盖交换机',
      type: 'Select',
      mode: 'multiple',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请输入覆盖交换机'
        }
      ],
      options: !options.device.loading ? (options.device.data || []).map(it => {
        return {
          label: it.name + '(' + it.fixed_asset_number + ')', value: it.fixed_asset_number
        };
      }) : []
    },
    {
      id: 'vlan',
      name: 'vlan',
      label: 'VLAN',
      type: 'TextInput',
      disabled: disabled || type === '_update',
      rules: [
        {
          required: true,
          message: '请输入VLAN'
        }
      ]
    },
    {
      id: 'version',
      name: 'version',
      label: 'IP版本',
      type: 'Radios',
      disabled: disabled || type === '_update',
      rules: [
        {
          required: true,
          message: '请选择IP版本'
        }
      ],
      radios: getSearchList(IP_VERSION)
    } 
  ];
  //查看详情多加几个字段显示
  if (disabled) {
    elements.unshift({
      id: 'idc_name',
      name: 'idc_name',
      label: '数据中心',
      type: 'TextInput',
      disabled: disabled
    });
    elements.push(
      {
        id: 'created_at',
        name: 'created_at',
        label: '创建时间',
        type: 'TextInput',
        disabled: disabled
      },
      {
        id: 'updated_at',
        name: 'updated_at',
        label: '更新时间',
        type: 'TextInput',
        disabled: disabled
      }
    );
  }
  return {
    name: 'form',
    id: 'form',
    elements: elements
  };
};
