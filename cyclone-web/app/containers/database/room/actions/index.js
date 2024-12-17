import _update from './create-update';
import _create from './create-update';
import _delete from './delete';
import _import from './import';
import detail from '../../common/room/detail';
import idc_detail from '../../common/idc/detail';
import accepted from './change-status';
import production from './change-status';
import abolished from './change-status';

export default {
  detail,
  idc_detail,
  _update,
  _delete,
  _create,
  _import,
  accepted,
  production,
  abolished
};
