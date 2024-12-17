import { post, get, getWithArgs } from 'common/xFetch2';
import { createTableAsyncAction } from 'utils/table-async-action';
import { notification } from 'antd';

async function getDetailInfo(state, action, dispatch) {
  try {
    dispatch({
      type: 'device/detail/load'
    });

    const res = await get(`/api/cloudboot/v1/devices/${action.payload}/combined`);
    if (res.status !== 'success') {
      return notification.error({ message: res.message });
    }
    dispatch({
      type: 'device/detail/load/success',
      payload: res.content
    });

  } catch (error) {
    console.log(error);
  }
}

export default {
  'device/detail/get': getDetailInfo
};
