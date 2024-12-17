import _update from './create-update';
import _create from './create-update';
import _delete from './delete';
import deletePort from './delete-port';
import _import from './import';
import _importPort from './import-port';
import detail from '../../common/usite/detail';
import cabinet_detail from '../../common/cabinet/detail';
import changeStatus from './change-status';
import remark from './remark';


export default {
  detail,
  cabinet_detail,
  _update,
  _delete,
  _create,
  _import,
  _importPort,
  deletePort,
  changeStatus,
  remark
};
