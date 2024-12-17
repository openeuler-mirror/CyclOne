import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createRegularReducer } from 'utils/regular-reducer';
import { createTableStore, createTableReducer } from 'utils/table-reducer';

const initialState = fromJS({
  loading: true,
  pendingTableData: createTableStore(),
  approvedTableData: createTableStore(),
  initiatedTableData: createTableStore(),
  approval_list: [
    {
      title: '数据中心',
      list: [
        {
          name: '数据中心裁撤',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/idc_abolish',
          permissionKey: 'button_approval_idc_abolish'
        },
        {
          name: '机房裁撤',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/server_room_abolish',
          permissionKey: 'button_approval_server_room_abolish'

        },
        {
          name: '机架关电',
          logo: 'assets/icon/guandian.png',
          link: 'approval/pages/cabinet_power_off',
          permissionKey: 'button_approval_cabinet_powerOff'

        },
        {
          name: '机架下线',
          logo: 'assets/icon/offline.png',
          link: 'approval/pages/cabinet_offline',
          permissionKey: 'button_approval_cabinet_offline'

        },
        {
          name: '网络区域下线',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/network_area_offline',
          permissionKey: 'button_approval_network_area_offline'

        }
      ]
    },
    {
      title: '物理机',
      list: [
        {
          name: '物理机搬迁',
          logo: 'assets/icon/move.png',
          link: 'approval/move',
          permissionKey: 'button_approval_physical_machine_move'

        },
        {
          name: '物理机退役',
          logo: 'assets/icon/tuiyi.png',
          link: 'approval/retire',
          permissionKey: 'button_approval_physical_machine_retirement'

        },
        {
          name: '物理机重装',
          logo: 'assets/icon/reinstall.png',
          link: 'device/entry?from=approval&type=reinstall',
          permissionKey: 'button_approval_physical_machine_reInstall'

        },
        {
          name: '物理机关电',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/device_power_off',
          permissionKey: 'button_approval_physical_machine_power_off'

        },
        {
          name: '物理机重启',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/device_restart',
          permissionKey: 'button_approval_physical_machine_restart'

        },
        {
          name: '回收待退役',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/device_recycle_pre_retire',
          permissionKey: 'button_approval_physical_machine_recycle_retire'

        },
        {
          name: '回收待搬迁',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/device_recycle_pre_move',
          permissionKey: 'button_approval_physical_machine_recycle_move'

        },
        {
          name: '回收重装',
          logo: 'assets/icon/reinstall.png',
          link: 'device/entry?from=approval&type=recycle',
          permissionKey: 'button_approval_physical_machine_recycle_reInstall'

        }
      ]
    },
    {
      title: '网络管理',
      list: [
        {
          name: 'IP回收',
          logo: 'assets/icon/reinstall.png',
          link: 'approval/pages/ip_unassign',
          permissionKey: 'button_ip_unassign'

        }
      ]
    }
  ]
});

const reducer = handleActions(
  {
    ...createTableReducer('approval/pending-table-data', 'pendingTableData'),
    ...createTableReducer('approval/approved-table-data', 'approvedTableData'),
    ...createTableReducer('approval/initiated-table-data', 'initiatedTableData')
  },
  initialState,
);

export default reducer;
