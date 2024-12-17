import React from 'react';
import echarts from 'echarts';

export default class TablePanel extends React.Component {
  constructor(props) {
    super(props);
    this.resultReport = null;
  }

  componentDidMount() {
    this.renderChart();
  }
  renderChart = () => {
    this.resultReport = echarts.init($(`#${this.props.id}`)[0]);
    $(window).resize(() => {
      const resultReportParent = $(`#${this.props.id}`);
      this.resultReport.resize(
        resultReportParent.height(),
        resultReportParent.width()
      );
    });
    // this.resultReport.showLoading();
    // const { inspections } = this.props;
    // if (!inspections.loading) {
    //   this.resultReport.hideLoading();
    //   this.resultReport.setOption(this.setReportOption(inspections.data));
    // }
    this.resultReport.setOption(this.setReportOption());
  };
  setReportOption = data => {
    const option = {
      color: [ '#4d9cff', '#a5afbc52' ],
      tooltip: {
        trigger: 'item',
        formatter: "{a} <br/>{b}: {c} ({d}%)"
      },
      legend: {
        show: false,
        x: 'right'
      },
      series: [
        {
          name: this.props.name,
          type: 'pie',
          radius: [ '48%', '70%' ],
          avoidLabelOverlap: false,
          label: {
            normal: {
              show: false,
              position: 'center'
            },
            emphasis: {
              show: false
            }
          },
          labelLine: {
            normal: {
              show: false
            }
          },
          data: [
            { value: 335, name: '已用容量' },
            { value: 210, name: '未用容量' }
          ]
        }
      ]
    };
    return option;
  };
  render() {
    return (
      <div className='chartWrapper'>
        <div id={this.props.id} style={{ height: 300 }} />
      </div>
    );
  }
}

