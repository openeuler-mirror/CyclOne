import React from 'react';
import { IDC_USAGE, IDC_STATUS, getSearchList } from 'common/enums';

export const formSchema = (disabled) => {
  let elements = [
    {
      id: 'name',
      name: 'name',
      label: '数据中心',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写数据中心名称'
        }
      ]
    },
    {
      id: 'usage',
      name: 'usage',
      label: '用途',
      type: 'Select',
      disabled: disabled,
      options: getSearchList(IDC_USAGE),
      rules: [
        {
          required: true,
          message: '请选择用途'
        }
      ]
    },
    {
      id: 'first_server_room',
      name: 'first_server_room',
      label: '一级机房',
      type: 'Select',
      mode: 'tags',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写一级机房'
        }
      ]
    },
    {
      id: 'vendor',
      name: 'vendor',
      label: '供应商',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写供应商'
        }
      ]
    }
  ];
  //查看详情多加几个字段显示
  if (disabled) {
    elements.push(
      {
        id: 'status',
        name: 'status',
        label: '状态',
        type: 'Select',
        disabled: disabled,
        options: getSearchList(IDC_STATUS)
      },
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
