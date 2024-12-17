import React from 'react';
import { Upload, Table, Button, Icon, Row, Col, notification, Alert } from 'antd';
import { post } from 'common/xFetch2';

/**
 * 适用于分三步导入的导入组件需求
 * importApi
 * uploadApi
 * previewApi
 * getColumns()
 */

export default class FileUpload extends React.Component {
  state = {
    fileList: [],
    tableList: [],
    fileName: '',
    disableUpload: false,
    errorMessage: null,
    importMessage: null
  };

  //最终提交
  handleUpload = () => {
    if (this.props.importApi) {
      console.log(this.props.importApi);
      console.log(this.state.fileName);
      post(this.props.importApi, { file_name: this.state.fileName }).then(res => {
        if (res.status !== 'success') {
          return notification.error({ message: res.message });
        }
        this.props.onSuccess();
        notification.success({ message: res.message });
      });
    } else {
      this.props.onSuccess(this.state.tableList);
    }
   
  };

  componentWillUnmount() {
    this.setState({
      fileList: [],
      tableList: [],
      fileName: '',
      disableUpload: false,
      errorMessage: null,
      importMessage: null
    });
  }

  render() {
    const props = {
      name: 'files[]',
      accept: '.xlsx',
      action: this.props.uploadApi,
      headers: {
        Authorization: localStorage.osinstallAuthAccessToken
      },
      onChange: (info) => {
        //清空
        this.setState({
          errorMessage: null,
          importMessage: null
        });

        let fileList = info.fileList;
        //只上传一个文件
        fileList = fileList.slice(-1);
        fileList = fileList.filter((file) => {
          if (file.response) {
            if (file.response.status === 'error') {
              return notification.error({ message: file.response.message });
            }
            //保存文件名称
            const fileName = file.response.content.result;
            this.setState({ fileName });

            //预览
            post(this.props.previewApi, { file_name: fileName, limit: 1000, offset: 0 }).then(res => {

              if (res.status === 'success') {
                //没有导入数据的时候会返回null
                this.setState({ disableUpload: false, tableList: res.content.content || [] });

                //导入的数据有误
                if (res.content.status === 'failure') {
                  this.setState({ disableUpload: true, errorMessage: res.content.message });
                } else {
                  this.setState({ disableUpload: false });
                  //导入结果显示
                  if (res.content.import_result) {
                    const importResult = res.content.import_result;
                    //暂时方案
                    if (typeof importResult !== 'object') {
                      return;
                    }
                    const importMessage = `可操作：${importResult.limit}，导入：${importResult.import_num}，新增：${importResult.import_num}，现存：${importResult.total_now}`;
                    this.setState({ importMessage: importMessage });
                  }
                }
              } else {
                //无法导入
                this.setState({ errorMessage: res.message, disableUpload: true });
              }

            });
            //过滤显示上传成功的文件
            return file.response.status === 'success';
          }
          return true;
        });
        this.setState({ fileList });
      }
    };

    const { errorMessage, importMessage, tableList } = this.state;
    let expandedRowKeys = [];
    if (tableList.length > 0) {
      tableList.map((item, index) => {
        if (item.content) {
          expandedRowKeys.push(index);
        }
      });
    }

    return (
      <div className='upload-form'>
        <Row>
          <Col span={2}>
            文件上传：
          </Col>
          <Col span={22}>
            <Upload {...props} fileList={this.state.fileList}>
              <Button>
                <Icon type='upload' /> 选择文件
              </Button>
            </Upload>
          </Col>
        </Row>

        <div style={{ marginBottom: 8 }}>
          {
            errorMessage && <Alert message={errorMessage} type='error' showIcon={true} closable={true}/>
          }
          {
            importMessage && <Alert message={importMessage} type='success' showIcon={true} closable={true}/>
          }
        </div>
        <Row>
          <Col span={24} className='no-wordbreak'>
            <Table
              dataSource={this.state.tableList}
              columns={this.props.getColumns()}
              pagination={{ showTotal: () => `共 ${this.state.tableList.length} 条` }}
              more={(record) => record.Content}
              expandedRowKeys={expandedRowKeys}
              expandedRowRender={record => <p style={{ margin: 0, color: '#ff3700' }} dangerouslySetInnerHTML={{ __html: record.content }} >
              </p>}
            />
          </Col>
        </Row>
        <Row>
          <Col span={24}>
            <div className='pull-right'>
              <Button onClick={() => this.props.onCancel()} style={{ marginRight: 8 }}>取消</Button>
              <Button type='primary' onClick={() => this.handleUpload()} disabled={this.state.tableList.length === 0 || this.state.disableUpload}>确定</Button>
            </div>
          </Col>
        </Row>
      </div>
    );
  }

}
