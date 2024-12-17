/**
 * 获取当前用户信息。
 *
 * Created by zhangrong on 16/12/25.
 */
import { base64decodeFunc } from 'common/base64';

export default function getCurrentUser() {
  const accout = base64decodeFunc(localStorage.OMC_ACCOUT);
  const name = base64decodeFunc(localStorage.CONT_NAME);
  const id = base64decodeFunc(localStorage.CONF_ID);

  return {
    accout, name, id
  };
}

