import React from 'react';
import { CAB_TYPE, CAB_STATUS, YES_NO, getSearchList } from 'common/enums';

export const formSchema = (disabled, options) => {
  let elements = [
    {
      id: 'network_area_id',
      name: 'network_area_id',
      label: '网络区域',
      type: 'Select',
      disabled: disabled,
      options: !options.network.loading ? options.network.data.map(it => {
        return {
          label: it.name + '（机房：' + it.server_room.name + '）', value: it.id
        };
      }) : [],
      rules: [
        {
          required: true,
          message: '请选择机房'
        }
      ]
    },
    {
      id: 'number',
      name: 'number',
      label: '机架编号',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写机架编号'
        }
      ]
    },
    {
      id: 'height',
      name: 'height',
      label: '机架高度',
      type: 'Number',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写机架高度'
        }
      ]
    },
    {
      id: 'type',
      name: 'type',
      label: '类型',
      type: 'Select',
      disabled: disabled,
      options: getSearchList(CAB_TYPE),
      rules: [
        {
          required: true,
          message: '请选择机房'
        }
      ]
    },
    {
      id: 'network_rate',
      name: 'network_rate',
      label: '网络速率',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写网络速率'
        }
      ]
    },
    {
      id: 'current',
      name: 'current',
      label: '电流',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写电流'
        }
      ]
    },
    {
      id: 'available_power',
      name: 'available_power',
      label: '可用功率',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写可用功率'
        }
      ]
    },
    {
      id: 'max_power',
      name: 'max_power',
      label: '峰值功率',
      type: 'TextInput',
      disabled: disabled,
      rules: [
        {
          required: true,
          message: '请填写峰值功率'
        }
      ]
    },
    {
      id: 'remark',
      name: 'remark',
      label: '备注',
      type: 'Textarea',
      disabled: disabled
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
        id: 'status',
        name: 'status',
        label: '机架状态',
        type: 'Select',
        disabled: disabled,
        options: getSearchList(CAB_STATUS)
      },
      {
        id: 'is_enabled',
        name: 'is_enabled',
        label: '是否启用',
        type: 'Select',
        disabled: disabled,
        options: getSearchList(YES_NO)
      },
      {
        id: 'enable_time',
        name: 'enable_time',
        label: '启用时间',
        type: 'TextInput',
        disabled: disabled
      },
      {
        id: 'is_powered',
        name: 'is_powered',
        label: '是否开电',
        type: 'Select',
        disabled: disabled,
        options: getSearchList(YES_NO)
      },
      {
        id: 'power_on_time',
        name: 'power_on_time',
        label: '开电时间',
        type: 'TextInput',
        disabled: disabled
      },
      {
        id: 'power_off_time',
        name: 'power_off_time',
        label: '关电时间',
        type: 'TextInput',
        disabled: disabled
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
