import { handleActions } from 'redux-actions';
import { fromJS } from 'immutable';
import { createRegularReducer } from 'utils/regular-reducer';

const initialState = fromJS({
  detailInfo: {
    loading: true,
    data: {
      device_page_resp: {

      },
      device_lifecycle_detail_page: {

      },      
      cpu: {
        physicals: []
      },
      disk: {

      },
      disk_slot: {
        items: []
      },
      memory: {
        items: []
      },
      nic: {
        items: []
      },
      oob: {
        network: [],
        user: [{}]
      },
      pci: {
        slots: []
      },
      hba: {
        item: []
      },
      raid: {
        items: []
      },
      motherboard: {

      }
    }
  },
  idc: {
    data: [],
    loading: true
  }
});

const reducer = handleActions(
  {
    ...createRegularReducer('device-detail/detail-info', 'detailInfo'),
    ...createRegularReducer('device-detail/idc', 'idc')
  },
  initialState,
);


export default reducer;
