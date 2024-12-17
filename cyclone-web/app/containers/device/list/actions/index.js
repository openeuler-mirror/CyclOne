import _delete from './delete';
import powerOn from '../../oob/actions/change-status';
import networkBoot from '../../oob/actions/change-status';
import reAccess from '../../oob/actions/change-status';
import powerOff from '../../oob/actions/change-status';
import _import from './import';
import pre_import from './pre_import';
import _importStoreRoom from './device_store_import';
import editStatus from './editStatus';
import editUsage from './editUsage';
import editHardwareRemark from './editHardwareRemark';


export default {
  powerOn,
  networkBoot,
  reAccess,
  powerOff, 
  _delete,
  _import,
  pre_import,
  _importStoreRoom,
  editStatus,
  editUsage,
  editHardwareRemark
};
