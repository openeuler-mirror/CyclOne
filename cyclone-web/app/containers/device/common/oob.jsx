import { Modal, notification } from 'antd';
import { get } from 'common/xFetch2';
import copy from 'copy-to-clipboard';

export default async function info(sn) {
  const res = await get(`/api/cloudboot/v1/devices/${sn}/oob-user`);
  if (res.status !== 'success') {
    return notification.error({ message: res.message });
  }
  Modal.info({
    title: '带外信息',
    okText: '知道了',
    content: (
      <div>
        <p>带外IP： <a href={`http://${res.content.ip}`} target='_blank'>{res.content.ip}</a></p>
        <p>带外用户： <span style={{ color: '#010f24' }}>{res.content.username}</span></p>
        <p>带外密码： <span style={{ color: '#010f24', marginRight: 8 }}>{res.content.password}</span>
          <a onClick={() => {
            if (copy(res.content.password)) {
              notification.success({ message: '复制成功' });
            } else {
              notification.error({ message: '复制失败' });
            }
          }}
          >复制密码</a>
        </p>
      </div>
    )
  });
}
