import React from 'react';
import { connect } from 'react-redux';
import Layout from 'components/layout/page-layout-tabs';
import { Tabs, Row, Col } from 'antd';
const TabPane = Tabs.TabPane;
import { Link } from 'react-router';
import SearchForm from './components/search-form';
import Table from './components/table';
import { getPermissonBtn } from 'common/utils';

class Container extends React.Component {

  componentDidMount() {
    // this.reload();
  }

  reload = () => {
    this.props.dispatch({
      type: 'approval/pending-table-data/get'
    });
    this.props.dispatch({
      type: 'approval/approved-table-data/get'
    });
    this.props.dispatch({
      type: 'approval/initiated-table-data/get'
    });
  };

  getContent = (data, type) => {
    return <div>
      <SearchForm type={type} dispatch={this.props.dispatch} initialValue={data.query}/>
      <Table
        dispatch={this.props.dispatch}
        tableData={data}
        type={type}
        reload={this.reload}
        userInfo={this.props.userInfo}
      />
    </div>;
  };
  renderTab = (key) => {
    this.props.dispatch({
      type: `approval/${key}-table-data/get`
    });
  };

  render() {
    const { approval_list, pendingTableData, approvedTableData, initiatedTableData } = this.props.data;
    return (
      <Layout title='审批管理'>
        <Tabs defaultActiveKey='apply' type='card' onTabClick={(key) => this.renderTab(key)}>
          <TabPane tab={<span><span className='tab_status tab_status_purple' />发起审批</span>} key='apply'>
            <div className='approval_list'>
              {approval_list.map(item => {
                return (
                  <Row gutter={8}>
                    <h3 className='approval_title'>{item.title}</h3>
                    {
                      item.list.map(it => {
                        if (getPermissonBtn(this.props.userInfo.permissions, it.permissionKey)) {
                          return <Col span='5'>
                            <Link to={it.link}>
                              <div className='home-card approval-card'>
                                <img src={it.logo} alt='' />
                                <span className='card-right text_single'>
                                  {it.name}
                                </span>
                              </div>
                            </Link>
                          </Col>;
                        } else {
                          return <Col span='5' title='无权限'>
                            <div className='home-card approval-card'>
                              <img src={it.logo} alt='' />
                              <span className='card-right text_single'>
                                {it.name}
                              </span>
                            </div>
                          </Col>;
                        }
                      })
                    }
                  </Row>
                );
              })}
            </div>
          </TabPane>
          <TabPane tab={<span><span className='tab_status tab_status_red' />待我审批</span>} key='pending'>
            {
              this.getContent(pendingTableData, 'pending')
            }
          </TabPane>
          <TabPane tab={<span><span className='tab_status tab_status_green' />我已审批</span>} key='approved'>
            {
              this.getContent(approvedTableData, 'approved')
            }
          </TabPane>
          <TabPane tab={<span><span className='tab_status tab_status_blue' />我发起的</span>} key='initiated'>
            {
              this.getContent(initiatedTableData, 'initiated')
            }
          </TabPane>
        </Tabs>
      </Layout>
    );
  }
}

function mapStateToProps(state) {
  return {
    data: state.get('approval').toJS(),
    userInfo: state.getIn([ 'global', 'userData' ]).toJS()
  };
}

function mapDispatchToProps(dispatch) {
  return {
    dispatch
  };
}

export default connect(mapStateToProps, mapDispatchToProps)(Container);
