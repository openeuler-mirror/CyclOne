/**
 * 简单封装下，然后在 index 中注册为 swal，和 ember 用法类似。
 *
 * Created by cxc on 17/6/7.
 */
import { notification } from 'antd';
export default function(type, message = '', description = '', duration) {
  if (duration !== null) {
    duration = duration ? duration : 2;
  }
  notification[type]({
    duration,
    message,
    description
  });
}
