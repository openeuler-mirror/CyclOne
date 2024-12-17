import React from 'react';
import { NET_STATUS, getSearchList } from 'common/enums';

export const formSchema = (disabled, options) => {
  let elements = [
    {
      id: 'name',
      name: 'name',
      label: '网络区域名称',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写网络区域名称'
        }
      ]
    },
    {
      id: 'server_room_id',
      name: 'server_room_id',
      label: '机房管理单元',
      type: 'Select',
      disabled: disabled,
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
      id: 'physical_area',
      name: 'physical_area',
      label: '关联物理区域',
      type: 'Select',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写关联物理区域'
        }
      ],
      mode: 'tags'
    },
    {
      id: 'status',
      name: 'status',
      label: '状态',
      type: 'Select',
      disabled: disabled,
      options: getSearchList(NET_STATUS),
      rules: [
        {
          required: true,
          message: '请选择状态'
        }
      ]
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
