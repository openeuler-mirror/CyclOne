import React from 'react';
import { renderFormDetail } from 'common/utils';
import JSONPretty from 'react-json-pretty';
import moment from 'moment';
import { TIME_FORMAT } from 'common/enums';

class MyForm extends React.Component {
  constructor(props) {
    super(props);
  }


  render() {
    const data = this.props.initialValue;

    return (
      <div>
        <h3 className='detail-title detail-title-drawer'>基本信息</h3>
        <div className='detail-info detail-info-24'>
          {renderFormDetail([
            {
              label: '接口地址',
              value: data.api
            },
            {
              label: '接口描述',
              value: data.description
            },
            {
              label: '请求方式',
              value: data.method
            },
            {
              label: '操作者',
              value: data.operator
            },
            {
              label: '操作时间',
              value: moment(data.created_at).format(TIME_FORMAT)
            },
            {
              label: '耗时（s）',
              value: data.time
            }
          ], 24)}
        </div>
        <h3 className='detail-title detail-title-drawer'>请求参数</h3>
        <div className='api-result'>
          <pre>
            { data.req_body && <JSONPretty json={JSON.parse(data.req_body)}></JSONPretty> }
          </pre>
        </div>

        <h3 className='detail-title detail-title-drawer'>结果详情</h3>
        <div className='detail-info detail-info-24'>
          {renderFormDetail([
            {
              label: '状态',
              value: data.status
            },
            {
              label: '消息',
              value: data.msg
            }
          ], 24)}
        </div>
        <div className='api-result'>
          <pre>
            <JSONPretty json={JSON.parse(data.result)}></JSONPretty>
          </pre>
        </div>
      </div>
    );
  }
}

export default MyForm;
