import React from 'react';
import { get } from 'common/xFetch2';
import { Card, Icon, Skeleton } from 'antd';
import { Link } from 'react-router';
import { getPermissonBtn } from 'common/utils';

class MyCard extends React.Component {

  render() {
    const { dataSource, loading } = this.props;
    return (
      <div className='hardware-card'>
        {
          dataSource.map(data => {
            return (
              <Card
                key={data.id}
                loading={loading}
                extra={(data.builtin === 'yes' || !getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_delete')) ? null : <Icon type='close' onClick={() => this.props.execAction('deleteTemplate', data)} />}
                title={data.name}
                style={{ width: 300 }}
                actions={data.builtin === 'yes' ?
                  [
                    <button disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_create')}
                      onClick={() => this.props.execAction('copyTemplate', data)}
                    ><Icon type='copy' />克隆</button>
                  ] :
                  [
                    <button disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_create')}
                      onClick={() => this.props.execAction('copyTemplate', data)}
                    ><Icon type='copy' theme='outlined' />克隆</button>,
                    <button disabled={!getPermissonBtn(this.props.userInfo.permissions, 'button_hardware_template_update')}
                      onClick={() => this.props.execAction('editTemplate', data)}
                    ><Icon type='edit' />编辑</button>
                  ]
                }
              >
                {/*<Skeleton loading={loading} active={true} paragraph={true}>*/}
                <Link to={`/template/hardware/detail/${data.id}`}>
                  {/*<p>厂商：{data.vendor}</p>*/}
                  {/*<p>型号：{data.model_name}</p>*/}
                  {/*<p>配置内容：{[...new Set((data.data || []).map(d => d.category))].join(' | ')}</p>*/}
                  <p>创建时间：{data.created_at}</p>
                  <p>修改时间：{data.updated_at}</p>
                </Link>
                {/*</Skeleton>*/}
              </Card>
            );
          })
        }
      </div>
    );
  }
}

export default MyCard;
