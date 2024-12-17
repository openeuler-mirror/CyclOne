import { get, post } from 'common/xFetch2';
import cookie from 'js-cookie';
import { notification } from 'antd';

/**
 * [getRedirectUrl description]
 * @return {[type]} [description]
 */
function getRedirectUrl(ssoUrl, customer) {
  const host = window.location.host;
  let loginHash = window.location.hash;
  loginHash = loginHash.replace('#/', '');
  const pathname = window.location.pathname;
  const protocol = window.location.protocol;

  const router = '/authenticated';
  const callbackUrl = encodeURIComponent(
    protocol + '//' + host + pathname + '#' + router + '?loginHash=' + loginHash
  );
  return ssoUrl + '?authCallbackUrl=' + callbackUrl + '&customer=' + customer;
}

/**
 * ssoLogin
 */
export function ssoLogin() {
  getSsoWebUrl().then(res => {
    const ssoUrl = res.content.url;
    const customer = res.content.customer;
    const url = getRedirectUrl(ssoUrl, customer);
    window.location.href = url;
  });
}

/**
 * [logout description]
 * @return {[type]} [description]
 */
export async function logout() {
  try {
    cookie.remove('access-token', '');
    document.cookie =
      'access-token=' + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
    delete localStorage.osinstallAuthAccessToken;
    ssoLogin();
  } catch (error) {
    notification.error({
      message: error.message
    });
  }
}

/**
 * 通过接口获取sso
 * @returns {*}
 */
export function getSsoWebUrl() {
  return get('/api/cloudboot/v1/system/login/settings');
}

/**
 * 获取菜单权限
 */
export function getMenuPermissions() {
  return get('/api/cloudboot/v1/users/info');
}
