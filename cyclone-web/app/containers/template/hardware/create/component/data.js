
export const category =
  [
    { category: 'RAID卡1', value: 'raid0' },
    { category: 'RAID卡2', value: 'raid1' },
    { category: 'RAID卡3', value: 'raid2' },
    { category: 'RAID卡4', value: 'raid3' }
    // { "category": "OOB", "value": "oob" },
    // { "category": "BIOS", "value": "bios" },
    // { "category": "FIRMWARE", "value": "firmware" },
    // { "category": "REBOOT", "value": "reboot" }
  ];
export const getTempData = () => {
  const raidAction = [
    { value: 'create_array', name: "阵列组" },
    { value: 'set_jbod', name: "直通磁盘列表" },
    { value: 'clear_settings', name: "擦除RAID配置" },
    { value: 'init_disk', name: "初始化逻辑磁盘" },
    { value: 'set_global_hotspare', name: "全局热备磁盘列表" }
  ];
  const RAIDList =
    {
      create_array: {
        action: "create_array",
        ci: "create_array",
        name: "阵列组",
        required: true,
        type: "list",
        children: [
          {
            ci: "level",
            name: "RAID级别",
            required: true,
            type: "select",
            options: [{ id: "raid0", name: "raid0" }, { id: "raid1", "name": "raid1" }, {
              id: "raid5",
              "name": "raid5"
            }, { id: "raid10", "name": "raid10" }, { id: "raid50", "name": "raid50" }]
          },
          { ci: "drives", name: "磁盘列表", required: true, type: "input" }
        ]
      },
      set_jbod: {
        action: "set_jbod",
        ci: "drives",
        name: "直通磁盘列表",
        desc: "直通磁盘名称列表，以逗号分隔",
        required: false,
        type: "input"
      },
      clear_settings:
      {
        action: "clear_settings",
        ci: "clear",
        name: "擦除RAID配置",
        desc: "是否先擦除raid配置",
        required: true,
        type: "radio"
      },
      init_disk: {
        action: "init_disk",
        ci: "init",
        name: "初始化逻辑磁盘",
        desc: "是否初始化逻辑磁盘",
        required: true,
        type: "radio"
      },
      set_global_hotspare: {
        action: "set_global_hotspare",
        ci: "drives",
        name: "全局热备磁盘列表",
        desc: "全局热备磁盘id列表，以逗号分隔",
        required: true,
        type: "input"
      }
    };
  return {
    raid0: {
      category: 'raid',
      block: 'RAID卡1',
      controller_index: '0',
      action: raidAction,
      metaData: RAIDList,
      initialValue: {
        category: 'raid',
        action: 'create_array',
        metadata: {
          controller_index: '0',
          level: 'raid0',
          drives: ''
        }
      }
    },
    raid1: {
      category: 'raid',
      block: 'RAID卡2',
      controller_index: '1',
      action: raidAction,
      metaData: RAIDList,
      initialValue: {
        category: 'raid',
        action: 'create_array',
        metadata: {
          controller_index: '1',
          level: 'raid0',
          drives: ''
        }
      }
    },
    raid2: {
      category: 'raid',
      block: 'RAID卡3',
      controller_index: '2',
      action: raidAction,
      metaData: RAIDList,
      initialValue: {
        category: 'raid',
        action: 'create_array',
        metadata: {
          controller_index: '2',
          level: 'raid0',
          drives: ''
        }
      }
    },
    raid3: {
      category: 'raid',
      block: 'RAID卡4',
      action: raidAction,
      controller_index: '3',
      metaData: RAIDList,
      initialValue: {
        category: 'raid',
        action: 'create_array',
        metadata: {
          controller_index: '3',
          level: 'raid0',
          drives: ''
        }
      }
    },
    oob: {
      category: 'oob',
      block: '配置项',
      action: [
        { name: '网络', value: "set_ip" },
        { name: '用户', value: "add_user" },
        { name: '冷重启', value: "reset_bmc" }
      ],
      metaData: {
        set_ip: {
          action: "set_ip", ci: "set_ip", name: "网络", required: true, type: "list",
          children: [
            { ci: "ip_src", name: "IP源", required: true, type: "select", options: [{ id: "dhcp", name: "dhcp" }, { id: "static", name: "静态" }] },
            { ci: "ip", name: "管理IP", required: true, type: "input", value: "<{manage_ip}>" },
            { ci: "netmask", name: "管理IP掩码", required: true, type: "input", value: "<{manage_netmask}>" },
            { ci: "gateway", name: "管理IP网关", required: true, type: "input", value: "<{manage_gateway}>" }
          ]
        },
        add_user: {
          action: "add_user", ci: "add_user", name: "用户", required: true, type: "list",
          children: [
            { ci: "username", name: "用户名", required: true, type: "input" },
            { ci: "password", name: "用户密码", required: true, type: "input" , inputType: 'password' },
            { ci: "privilege_level", name: "用户权限级别", required: true, type: "select", options: [{ id: "1", name: "1(callback level)" }, { id: "2", name: "2(user level)" }, { id: "3", name: "3(operator level)" }, { id: "4", name: "4(admin level)" }, { id: "5", name: "5(OEM proprietary level)" }] }
          ]
        },
        reset_bmc: { action: "reset_bmc", ci: "reset", name: "冷重启", desc: "冷重启", required: true, type: "radio" }
      },
      initialValue: {
        category: 'oob',
        action: 'set_ip',
        metadata: {
          gateway: '<{manage_gateway}>',
          ip_src: 'static',
          ip: '<{manage_ip}>',
          netmask: '<{manage_netmask}>'
        }
      }
    },
    bios: {
      category: 'bios',
      block: {
        type: 'select',
        name: 'Custom',
        defaultValue: 'Custom',
        options: [
          { name: 'Custom', value: 'custom' },
          { name: 'Dell', value: 'Dell' },
          { name: 'HP', value: 'HP' },
          { name: 'H3C', value: 'H3C' }
        ]
      },
      action: [{ name: '', value: '' }],
      metaData: {
        'custom': [

        ],
        "Dell": [
          {
            "ci": "ProcVirtualization",
            "name": "虚拟化",
            "desc": "是否开启虚拟化",
            "type": "radio"
          },
          {
            "ci": "ProcCores",
            "name": "每个处理器核心数",
            "desc": "每个处理器的核心数量",
            "type": "select",
            "options": [
              { "id": "1", "name": "1" },
              { "id": "ALL", "name": "全部" }
            ]
          },
          {
            "ci": "ProcX2Apic",
            "name": "CPU中断增强",
            "type": "radio"
          },
          {
            "ci": "AcPwrRcvry",
            "name": "交流电源恢复(AcPwrRcvry)",
            "type": "select",
            "options": [
              { "id": "last", "name": "上一次(last)" },
              { "id": "on", "name": "开启(on)" },
              { "id": "off", "name": "关闭(off)" }
            ]
          },
          {
            "ci": "PowerSaver",
            "name": "省电模式(PowerSaver)",
            "type": "radio"
          },
          {
            "ci": "SysProfile",
            "name": "系统配置文件设置(SysProfile)",
            "type": "select",
            "options": [
              { "id": "custom", "name": "自定义(custom)" },
              { "id": "PerfOptimized", "name": "最大性能(PerfOptimized)" }
            ]
          },
          {
            "ci": "ProcCStates",
            "name": "处理器在C电源状态下运行(ProcCStates)",
            "type": "radio"
          },
          {
            "ci": "ProcPwrPerf",
            "name": "处理器电源性能(ProcPwrPerf)",
            "type": "select",
            "options": [
              { "id": "sysdbpm", "name": "基于系统的电源管理(sysdbpm)" },
              { "id": "osdbpm", "name": "基于OS的电源管理(osdbpm)" },
              { "id": "maxperf", "name": "最大性能的电源管理(maxperf)" },
              { "id": "hwpdbpm", "name": "基于硬件的电源管理(hwpdbpm)" }
            ]
          },
          {
            "ci": "EnergyPerformanceBias",
            "name": "电源性能管理(EnergyPerformanceBias)",
            "type": "select",
            "options": [
              { "id": "maxpower", "name": "最大功率(maxpower)" },
              { "id": "balancedperformance", "name": "性能平衡(balancedperformance)" },
              { "id": "balancedefficiency", "name": "功耗平衡(balancedefficiency)" },
              { "id": "lowpower", "name": "低功耗(lowpower)" }
            ]
          },
          {
            "ci": "SerialComm",
            "name": "串行通信(SerialComm)",
            "type": "select",
            "options": [
              { "id": "onnoconredir", "name": "开启所有串口通信(onnoconredir)" },
              { "id": "onconredirauto", "name": "自识别串口通信(onconredirauto)" },
              { "id": "onconredircom1", "name": "开启com1端口(onconredircom1)" },
              { "id": "onconredircom2", "name": "开启com2端口(onconredircom2)" },
              { "id": "off", "name": "关闭(off)" }
            ]
          },
          {
            "ci": "WriteCache",
            "name": "磁盘写缓存(WriteCache)",
            "type": "radio"
          }
        ],
        "HP": [
          {
            "ci": "CPU_Virtualization",
            "name": "CPU虚拟化(CPU_Virtualization)",
            "type": "radio"
          },
          {
            "ci": "Intel_VT-d2",
            "name": "I/O虚拟化(Intel_VT-d2)",
            "type": "radio"
          },
          {
            "ci": "Intel_Hyperthreading",
            "name": "超线程(Intel_Hyperthreading)",
            "type": "radio"
          },
          {
            "ci": "BIOS_Console",
            "name": "BIOS控制台(BIOS_Console)",
            "type": "select",
            "options": [
              { "id": "Disabled", "name": "关闭(Disabled)" },
              { "id": "Auto", "name": "自适应(Auto)" }
            ]
          },
          {
            "ci": "HP_Power_Regulator",
            "name": "HP 功率调节器(HP_Power_Regulator)",
            "type": "select",
            "options": [
              { "id": "HP_Static_High_Performance_Mode", "name": "HP 静态高性能模式(HP_Static_High_Performance_Mode)" }
            ]
          },
          {
            "ci": "Redundant_Power_Supply_Mode",
            "name": "冗余电源模式(Redundant_Power_Supply_Mode)",
            "type": "select",
            "options": [
              { "id": "High_Efficiency_Even_Standby", "name": "高效模式(偶数电源待机)(High_Efficiency_Even_Standby)" },
              { "id": "Balanced_Mode", "name": "平衡模式(Balanced_Mode)" }
            ]
          },
          {
            "ci": "HW_Prefetch",
            "name": "硬件预取(HW_Prefetch)",
            "type": "radio"
          },
          {
            "ci": "Adjacent_Sector_Prefetch",
            "name": "相邻扇区预取(Adjacent_Sector_Prefetch)",
            "type": "radio"
          },
          {
            "ci": "Intel_Processor_Core_Disable",
            "name": "禁用增强的处理器内核(Intel_Processor_Core_Disable)",
            "type": "select",
            "options": [
              { "id": "All_Cores_Enabled", "name": "开启所有处理器(All_Cores_Enabled)" }
            ]
          },
          {
            "ci": "Intel_Processor_Turbo_Mode",
            "name": "Intel 睿频加速模式(Intel_Processor_Turbo_Mode)",
            "type": "radio"
          },
          {
            "ci": "HP_Power_Profile",
            "name": "HP 电源配置文件(HP_Power_Profile)",
            "type": "select",
            "options": [
              { "id": "Maximum_Performance", "name": "最高性能(Maximum_Performance)" }
            ]
          },
          {
            "ci": "Intel_Minimum_Processor_Idle_Power_State",
            "name": "最低处理器空闲电源状态(Intel_Minimum_Processor_Idle_Power_State)",
            "type": "select",
            "options": [
              { "id": "No_C-States", "name": "无C 状态(No_C-States)" }
            ]
          },
          {
            "ci": "Intel_Turbo_Boost_Optimization",
            "name": "Intel 睿频加速优化(Intel_Turbo_Boost_Optimization)",
            "type": "select",
            "options": [
              { "id": "Optimized_for_Performance", "name": "为性能优za化(Optimized_for_Performance)" }
            ]
          }
        ],
        "H3C": [
          {
            "ci": "CPU_Virtualization",
            "name": "CPU虚拟化(CPU_Virtualization)",
            "type": "radio"
          },
          {
            "ci": "Intel_VT-d2",
            "name": "I/O虚拟化(Intel_VT-d2)",
            "type": "radio"
          },
          {
            "ci": "Intel_Hyperthreading",
            "name": "超线程(Intel_Hyperthreading)",
            "type": "radio"
          },
          {
            "ci": "BIOS_Console",
            "name": "BIOS控制台(BIOS_Console)",
            "type": "select",
            "options": [
              { "id": "Disabled", "name": "关闭(Disabled)" },
              { "id": "Auto", "name": "自适应(Auto)" }
            ]
          },
          {
            "ci": "HP_Power_Regulator",
            "name": "HP 功率调节器(HP_Power_Regulator)",
            "type": "select",
            "options": [
              { "id": "HP_Static_High_Performance_Mode", "name": "HP 静态高性能模式(HP_Static_High_Performance_Mode)" }
            ]
          },
          {
            "ci": "Redundant_Power_Supply_Mode",
            "name": "冗余电源模式(Redundant_Power_Supply_Mode)",
            "type": "select",
            "options": [
              { "id": "High_Efficiency_Even_Standby", "name": "高效模式(偶数电源待机)(High_Efficiency_Even_Standby)" },
              { "id": "Balanced_Mode", "name": "平衡模式(Balanced_Mode)" }
            ]
          },
          {
            "ci": "HW_Prefetch",
            "name": "硬件预取(HW_Prefetch)",
            "type": "radio"
          },
          {
            "ci": "Adjacent_Sector_Prefetch",
            "name": "相邻扇区预取(Adjacent_Sector_Prefetch)",
            "type": "radio"
          },
          {
            "ci": "Intel_Processor_Core_Disable",
            "name": "禁用增强的处理器内核(Intel_Processor_Core_Disable)",
            "type": "select",
            "options": [
              { "id": "All_Cores_Enabled", "name": "开启所有处理器(All_Cores_Enabled)" }
            ]
          },
          {
            "ci": "Intel_Processor_Turbo_Mode",
            "name": "Intel 睿频加速模式(Intel_Processor_Turbo_Mode)",
            "type": "radio"
          },
          {
            "ci": "HP_Power_Profile",
            "name": "HP 电源配置文件(HP_Power_Profile)",
            "type": "select",
            "options": [
              { "id": "Maximum_Performance", "name": "最高性能(Maximum_Performance)" }
            ]
          },
          {
            "ci": "Intel_Minimum_Processor_Idle_Power_State",
            "name": "最低处理器空闲电源状态(Intel_Minimum_Processor_Idle_Power_State)",
            "type": "select",
            "options": [
              { "id": "No_C-States", "name": "无C 状态(No_C-States)" }
            ]
          },
          {
            "ci": "Intel_Turbo_Boost_Optimization",
            "name": "Intel 睿频加速优化(Intel_Turbo_Boost_Optimization)",
            "type": "select",
            "options": [
              { "id": "Optimized_for_Performance", "name": "为性能优化(Optimized_for_Performance)" }
            ]
          }
        ]

      },
      initialValue: {
        category: 'bios',
        action: 'set_bios',
        metadata: {
          manufacturer: 'custom',
          custom: 'YES/NO', //老系统字段
          key: '',
          value: ''
        }
      }
    },
    firmware: {
      category: 'firmware',
      block: '',
      action: [{ name: '', value: '' }],
      metaData: {},
      initialValue: {
        category: 'firmware',
        action: 'update_package',
        metadata: {
          file: '',
          category: '',
          category_desc: '',
          expected: ''
        }
      }
    },
    reboot: {
      category: 'reboot',
      block: 'reboot',
      action: [{ name: '', value: '' }],
      metaData: {},
      initialValue: {

      }
    }
  };
};

export const getFormData = (cards, values) => {
  const data = [];
  cards.map((card, index) => {
    switch (card.category) {
    case 'raid':
      const action = values[`raid-action-${card.uuid}`];
      switch (action) {
      case 'create_array':
        data[index] = {
          category: 'raid',
          action: action,
          metadata: {
            controller_index: card.controller_index,
            level: values[`raid-create_array-level-${card.uuid}`],
            drives: values[`raid-create_array-drives-${card.uuid}`]
          }
        };
        break;
      case 'set_jbod':
        data[index] = {
          category: 'raid',
          action: action,
          metadata: {
            controller_index: card.controller_index,
            drives: values[`raid-set_jbod-drives-${card.uuid}`]
          }
        };
        break;
      case 'clear_settings':
        data[index] = {
          category: 'raid',
          action: action,
          metadata: {
            controller_index: card.controller_index,
            clear: values[`raid-clear_settings-clear-${card.uuid}`] ? 'ON' : 'OFF'
          }
        };
        break;
      case 'init_disk':
        data[index] = {
          category: 'raid',
          action: action,
          metadata: {
            controller_index: card.controller_index,
            init: values[`raid-init_disk-init-${card.uuid}`] ? 'ON' : 'OFF'
          }
        };
        break;
      case 'set_global_hotspare':
        data[index] = {
          category: 'raid',
          action: action,
          metadata: {
            controller_index: card.controller_index,
            drives: values[`raid-set_global_hotspare-drives-${card.uuid}`]
          }
        };
        break;
      default:
      }
      break;
    case 'oob':
      const action_oob = values[`oob-action-${card.uuid}`];
      switch (action_oob) {
      case 'set_ip':
        data[index] = {
          category: 'oob',
          action: action_oob,
          metadata: {
            gateway: values[`oob-set_ip-gateway-${card.uuid}`],
            ip_src: values[`oob-set_ip-ip_src-${card.uuid}`],
            ip: values[`oob-set_ip-ip-${card.uuid}`],
            netmask: values[`oob-set_ip-netmask-${card.uuid}`]
          }
        };
        break;
      case 'add_user':
        data[index] = {
          category: 'oob',
          action: action_oob,
          metadata: {
            privilege_level: values[`oob-add_user-privilege_level-${card.uuid}`],
            username: values[`oob-add_user-username-${card.uuid}`],
            password: values[`oob-add_user-password-${card.uuid}`]
          }
        };
        break;
      case 'reset_bmc':
        data[index] = {
          category: 'oob',
          action: action_oob,
          metadata: {
            category: 'cold',
            reset: values[`oob-reset_bmc-reset-${card.uuid}`] ? 'ON' : 'OFF'
          }
        };
        break;
      default:
      }
      break;
    case 'bios':
      let key = '';
      const custom = values[`bios-custom-${card.uuid}`];
      if (custom === 'custom') {
        key = values[`bios-value-${card.uuid}`];
      } else if (custom === 'Dell') {
        key = values[`bios-value-${card.uuid}`] ? 'enable' : 'disable';
      } else {
        key = values[`bios-value-${card.uuid}`] ? 'Enabled' : 'Disabled';
      }
      data[index] = {
        category: 'bios',
        action: 'set_bios',
        metadata: {
          manufacturer: custom,
          key: values[`bios-key-${card.uuid}`],
          value: key
        }
      };
      break;
    case 'firmware':
      const category = JSON.parse(values[`firmware-category-${card.uuid}`] || '{}');
      data[index] = {
        category: 'firmware',
        action: 'update_package',
        metadata: {
          file: values[`firmware-file-${card.uuid}`],
          category: category.name,
          category_desc: category.value,
          expected: values[`firmware-expected-${card.uuid}`]
        }
      };
      break;
    case 'reboot':
      data[index] = {
        category: 'reboot',
        action: 'reboot',
        metadata: {
        }
      };
      break;
    default:
    }
  });
  return data;
};


