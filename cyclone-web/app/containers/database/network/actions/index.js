import _update from './create-update';
import _create from './create-update';
import _delete from './delete';
import _import from './import';
import detail from '../../common/network/detail';
import room_detail from '../../common/room/detail';
import production from './change-status';
import offline from './change-status';

export default {
  detail,
  room_detail,
  _update,
  _delete,
  _create,
  _import,
  production,
  offline
};
