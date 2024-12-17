import React from 'react';
import FormGenerator from './FormGenerator';
import { notification } from 'antd';
import { post, put, del } from './common/xFetch2';
const nil = () => {};
export default class CustomForm extends React.Component {

  constructor(props) {
    super(props);
    if (this.props.schema) {
      this.props.schema.forceUpdate = this.forceUpdate.bind(this);
    }
    this.state = {
      loading: false,
      schema: this.props.schema,
      message: this.props.message
    };
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.schema) {
      if (this.state.schema !== nextProps.schema) {
        nextProps.schema.forceUpdate = this.forceUpdate.bind(this);
        this.setState({
          schema: nextProps.schema
        });
      }
    }
  }
  handleTest = values => {

    if (this.state.testLoading) {
      return;
    }

    this.setState({
      testLoading: true
    });

    let api = this.props.testApi;
    let method = post;

    values = values || {};
    Object.keys(values).forEach(k => {
      const v = values[k];
      if (v && v.toDate) {
        values[k] = v.format('YY-MM-DD hh:mm:ss');
      }
    });

    const onError = ret => {
      this.setState({
        testLoading: false
      });
      notification.error({
        message: ret.message
      });
    };

    const data = {
      ...values,
      ...this.props.extraData
    };
    method(api, data).then(ret => {
      if (ret.status === 'success') {
        this.setState({
          testLoading: false
        });
        notification.success({
          message: '执行成功'
        });
      } else {
        onError(ret);
      }
    }, onError
    );
  }
  handleSubmit = values => {
    if (this.props.beforeSubmit) {
      const flag = this.props.beforeSubmit(values);
      if (flag) {
        return;
      }
    }
    if (this.state.loading) {
      return;
    }

    this.setState({
      loading: true,
      message: this.props.loadingMessage || this.props.message
    });


    let api = this.props.api;

    let method = post;
    if (this.props.method === 'put') {
      method = put;
    } else if (this.props.method === 'delete') {
      method = del;
    }

    values = values || {};
    Object.keys(values).forEach(k => {
      const v = values[k];
      if (v && v.toDate) {
        values[k] = v.format('YY-MM-DD hh:mm:ss');
      }
      //去除空格
      if (this.props.trim) {
        if (typeof v === 'string') {
          values[k] = v.trim();
        } else if (typeof v === 'object' && Array.isArray(v)) {
          values[k] = v.map(it => it.trim());
        }
      }
    });



    if (this.props.getApi) {
      api = this.props.getApi({
        ...values,
        ...this.props.extraData
      });
    }

    const onError = ret => {
      this.setState({
        loading: false,
        message: this.props.errorMessage || this.props.message
      });

      if (this.props.beforeError) {
        this.props.beforeError(ret);
      } else {
        notification.error({
          message: ret.message
        });
      }
      //返回失败时不关闭弹窗，保留数据
      //this.props.onError(ret);
    };

    let data = null;
    //处理数组格式的数据
    if (this.props.isArray) {
      data = this.props.dataArray;
    } else if (this.props.nullData) {
      //处理{}
      data = {};
    } else {
      //其他
      data = {
        ...values,
        ...this.props.extraData
      };
    }
    method(api, data).then(ret => {
      if (ret.status === 'success') {
        if (this.props.itemResult && ret.item.success === false) {
          return onError(ret);
        }
        this.setState({
          loading: false,
          message: this.props.successMessage || this.props.message
        });

        // 是否要拦截成功的处理逻辑
        let isProcessSuccess = true;
        if (this.props.beforeSuccess) {
          isProcessSuccess = this.props.beforeSuccess(ret);
        }

        if (isProcessSuccess) {
          notification.success({
            message: '执行成功'
          });
          this.props.onSuccess(ret);
        }
      } else {
        onError(ret);
      }
    }, onError
    );
  };

  render() {
    const props = this.props;
    let $content = null;
    //返回失败时不关闭弹窗，保留数据
    if (props.type === 'form') {
      $content = this.renderForm();
    }

    if (props.type === 'confirm') {
      $content = this.renderConfirm();
    }

    return <div className='custom-form'>{$content}</div>;
  }

  renderForm = () => {
    const props = this.props;
    const initialValue = props.initialValue;
    const schema = props.schema;
    //提示语补充
    if (schema.elements && schema.elements.length > 0) {
      schema.elements.map(ele => {
        if (ele.rules && ele.rules.length > 0) {
          ele.rules.map(rule => {
            if (rule.required && !rule.message) {
              rule.message = `${ele.label}不能为空`;
            }
          });
        }
      });
    }

    return (
      <div className='form'>
        <div className='info'>{props.message}</div>
        <FormGenerator
          schema={schema}
          initialValue={initialValue}
          onSubmit={this.handleSubmit}
          onTest={this.handleTest}
          showTest={this.props.showTest}
        />
      </div>
    );
  };


  renderConfirm = () => {
    const props = this.props;
    return (
      <div className='ant-confirm ant-confirm-confirm'>
        <div className='ant-modal-body' style={{ padding: 12 }}>
          <div className='ant-confirm-body-wrapper'>
            <div className='ant-confirm-body'>
              <i className='anticon anticon-question-circle' />
              <span className='ant-confirm-title'>
                {this.state.message}
              </span>
              <div className='ant-confirm-content'>
                {props.extraMessage}</div>
            </div>
            <div className='ant-confirm-btns'>
              <button type='button' className='ant-btn' disabled={this.state.loading} onClick={ev => {
                this.props.onCancel();
              }}
              ><span>取 消</span></button>
              <button type='button' className='ant-btn ant-btn-primary'
                disabled={this.state.loading} onClick={ev => {
                  this.handleSubmit();
                }}
              ><span>确 定</span></button>
            </div>
          </div>
        </div>
      </div>
    );
  };
}

CustomForm.defaultProps = {
  reload: nil,
  handleSubmit: nil,
  successMessage: '执行成功',
  triggerClassName: '',
  type: 'form',
  extraData: {},
  dataArray: []
};
