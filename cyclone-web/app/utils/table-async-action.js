import { notification } from 'antd';
import { post, get, getWithArgs } from 'common/xFetch2';
import { fromJS } from 'immutable';

const pause = time => {
  return new Promise(resolve => {
    setTimeout(ev => {
      resolve();
    }, time);
  });
};

export function createTableAsyncAction(options) {
  const {
    actionNamePrefix = 'table',
    tableDataPath = ['tableData'], // ["project", "tableData"]
    datasource,
    noMoreQuery,
    getDatasource,
    pageParameter,
    fetchMethod = getWithArgs,
    getExtraQuery = (state, action) => {
      return {};
    } // 在查询 table 数据的时候需要获取额外的数据
  } = options;

  /**
   * 加载 Table 数据
   */
  async function getTableData(state, action, dispatch, select, url) {
    try {
      dispatch({
        type: `${actionNamePrefix}/load`
      });

      state = select();
      const extraQuery = getExtraQuery(state, action);
      const pagination = state.getIn([ ...tableDataPath, 'pagination' ]).toJS();
      let query = state.getIn([ ...tableDataPath, 'query' ]);
      let sorter = state.getIn([ ...tableDataPath, 'sorter' ]);

      if (query) {
        query = query.toJS();
      } else {
        query = {};
      }

      // 获取数据方法
      // 接口可能不支持传入limit, page 等额外属性，因此需要加个判断，虽然丑了点
      const fetch = options.fetchMethod || getWithArgs;

      let realDatasource = datasource || url;
      if (typeof getDatasource === 'function') {
        realDatasource = getDatasource(state, action);
      }
      //boot新的接口分页参数传page，page_size，旧的接口传Limit,Offset,Offset从0开始
      let pageSizeKey = 'page_size';
      let pageKey = 'page';
      let page = pagination.page;
      if (pageParameter) {
        pageKey = pageParameter.page;
        pageSizeKey = pageParameter.pageSize;
        page = page - 1;
      }

      const ret = await fetch(
        realDatasource,
        noMoreQuery
          ? { ...extraQuery, ...query }
          : {
            ...query,
            ...extraQuery,
            ...sorter,
            [pageKey]: page,
            [pageSizeKey]: pagination.pageSize
          }
      );
      const content = ret.content || ret.Content;
      dispatch({
        type: `${actionNamePrefix}/load/success`,
        payload: {
          content: content.records || content.list,
          pagination: {
            pageSize: content.page_size || pagination.pageSize,
            page: content.page || pagination.page,
            total: content.total_records || content.recordCount
          }
        }
      });

      if (action.payload && action.payload.cb) {
        action.payload.cb(ret);
      }
    } catch (err) {
      notification.error({
        message: err.message
      });
      console.log(err);
    }
  }

  /**
   * 搜索
   */
  async function searchTable(state, action, dispatch, select, url) {
    try {
      dispatch([
        {
          type: `${actionNamePrefix}/set-query`,
          payload: action.payload
        },
        {
          type: `${actionNamePrefix}/set-page`,
          payload: {
            page: 1
          }
        }
      ]);
      await pause(100);
      await getTableData(state, action, dispatch, select, url);
    } catch (err) {
      notification.error({
        message: err.message
      });
      console.log(err);
    }
  }

  /**
   * 操作定义分页
   */
  async function changePage(state, action, dispatch, select, url) {
    try {
      dispatch({
        type: `${actionNamePrefix}/set-page`,
        payload: action.payload
      });
      await pause(100);
      await getTableData(state, action, dispatch, select, url);
    } catch (err) {
      notification.error({
        message: err.message
      });
      console.log(err);
    }
  }

  /**
   * 操作定义分页大小
   */
  async function changePageSize(state, action, dispatch, select, url) {
    try {
      dispatch({
        type: `${actionNamePrefix}/set-page-size`,
        payload: action.payload
      });
      await pause(100);
      await getTableData(state, action, dispatch, select, url);
    } catch (err) {
      notification.error({
        message: err.message
      });
      console.log(err);
    }
  }

  return {
    [`${actionNamePrefix}/change-page-size`]: changePageSize,
    [`${actionNamePrefix}/change-page`]: changePage,
    [`${actionNamePrefix}/search`]: searchTable,
    [`${actionNamePrefix}/get`]: getTableData
  };
}
