import React from 'react';
import echarts from 'echarts';

export default class TablePanel extends React.Component {
  constructor(props) {
    super(props);
    this.resultReport = null;
  }

  componentDidUpdate() {
    this.renderChart();
  }
  renderChart = () => {
    this.resultReport = echarts.init($('#resultReport')[0]);
    $(window).resize(() => {
      const resultReportParent = $('#resultReport');
      this.resultReport.resize(
        resultReportParent.height(),
        resultReportParent.width()
      );
    });
    this.resultReport.showLoading();
    const { inspections } = this.props;
    if (!inspections.loading) {
      this.resultReport.hideLoading();
      this.resultReport.setOption(this.setReportOption(inspections.data));
    }
  };
  setReportOption = data => {
    const option = {
      color: [ 'rgba(139,212,109,1)', 'rgba(255,208,77,1)', 'rgba(243,111,99,0.9)' ],
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          label: {
            backgroundColor: '#6a7985'
          }
        }
      },
      legend: {
        data: [ '正常', '警告', '异常' ],
        right: '3%',
        top: '2%'
      },
      grid: {
        left: '3%',
        right: '3%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: data.map(item => item.date)
      },
      yAxis: [
        {
          type: 'value'
        }
      ],
      series: [
        {
          name: '正常',
          type: 'line',
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: 'rgba(139,212,109,0.1)' // 0% 处的颜色
              }, {
                offset: 1, color: 'rgba(139,212,109,0)' // 100% 处的颜色
              }],
              globalCoord: false // 缺省为 false
            }
          },
          data: data.map(it => it.nominal_count)
        },
        {
          name: '警告',
          type: 'line',
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: 'rgba(255,208,77,0.1)' // 0% 处的颜色
              }, {
                offset: 1, color: 'rgba(255,208,77,0)' // 100% 处的颜色
              }],
              globalCoord: false // 缺省为 false
            }
          },
          data: data.map(it => it.warning_count)
        },
        {
          name: '异常',
          type: 'line',
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [{
                offset: 0, color: 'rgba(243,111,99,0.1)' // 0% 处的颜色
              }, {
                offset: 1, color: 'rgba(243,111,99,0)' // 100% 处的颜色
              }],
              globalCoord: false // 缺省为 false
            }
          },
          data: data.map(it => it.critical_count)
        }
      ]
    };
    return option;
  };
  render() {
    return (
      <div className='chartWrapper'>
        <div id='resultReport' style={{ height: 350 }} />
      </div>
    );
  }
}

