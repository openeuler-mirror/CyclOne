import _update from './create-update';
import _create from './create-update';
import _delete from './delete';
import _import from './import';
import detail from '../../common/cabinet/detail';
import network_detail from '../../common/network/detail';
import offline from './change-status';
import enabled from './change-status';
import locked from './change-status';
import powerOn from './powerOn';
import powerOff from './powerOff';
import changeType from './change-type';
import remark from './remark';

export default {
  detail,
  _update,
  _delete,
  _create,
  _import,
  enabled,
  offline,
  locked,
  powerOn,
  powerOff,
  network_detail,
  changeType,
  remark
};
