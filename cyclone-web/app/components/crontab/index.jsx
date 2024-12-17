import React from 'react';
import { Tabs, Popover, Form, Radio, Icon, Input, Select, Button, Checkbox } from 'antd';
const TabPane = Tabs.TabPane;
import { post, put, get } from 'common/xFetch2';
const FormItem = Form.Item;
const RadioGroup = Radio.Group;
const Option = Select.Option;
import styles from './styles.less';
const tableFormItemLayout = {
  labelCol: {
    xs: { span: 3 }
  },
  wrapperCol: {
    xs: { span: 19 }
  }
};
const radioStyle = {
  display: 'block',
  height: '36px'
};

/**
 *
 *  <Crontab
 form={this.props.form}
 initialValue={{crontab:string,crontabUi:JSONstring}}
 handleClick={(values)=>this.setState({crontabUi: values})}
 />
 *
 */

export default class Crontab extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      crontabExpression: '0 1/5 * * * ?',
      'type': 'define', // define：指定，custom：自定义，下面信息金在define时存在
      'minute': {
        'type': 'loop', // loop: 循环，define：定义
        'data': {
          'from': '1', // 仅在type为loop时存在
          'to': '5', // 仅在type为loop时存在
          'index': [] // 仅在type为define时使用，为勾选的所有的数据
        }
      },
      'hour': {
        'type': 'loop',
        'data': {
          'index': [] // 仅在type为define时使用，为勾选的所有的数据
        }
      },
      'day': {
        'type': 'loop',
        'data': {
          'index': [] // 仅在type为define时使用，为勾选的所有的数据
        }
      },
      'month': {
        'type': 'loop',
        'data': {
          'index': [] // 仅在type为define时使用，为勾选的所有的数据
        }
      },
      'weekly': {
        'type': 'unused', //unused表示不使用星期
        'data': {
          'index': [] // 仅在type为define时使用，为勾选的所有的数据
        }
      }
    };
    this.props.handleClick(this.state);
  }
  componentDidMount() {
    if (this.props.initialValue && this.props.initialValue.crontabUi && typeof this.props.initialValue.crontabUi !== 'string') {
      const data = JSON.parse(this.props.initialValue.crontabUi);
      if (typeof data === String) {
        return;
      }
      this.state = data;
    }
  }

  onMinTypeChange = e => {
    const minute = this.state.minute;

    //这一块处理自动生成表达式
    const minuteData = {
      minute: {
        ...minute,
        type: e.target.value
      }
    };
    const data = {
      ...this.state,
      ...minuteData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  };
  onHourTypeChange = e => {
    const hour = this.state.hour;

    //这一块处理自动生成表达式
    const hourData = {
      hour: {
        ...hour,
        type: e.target.value
      }
    };
    const data = {
      ...this.state,
      ...hourData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  onDayTypeChange = e => {
    const day = this.state.day;

    //校验
    const weekly = this.state.weekly;
    let weeklyData = { weekly: { ...weekly } };
    if (e.target.value !== 'unused') {
      weeklyData = {
        weekly: {
          ...weekly,
          type: 'unused'
        }
      };
    } else {
      let type = 'loop';
      if (weekly.data.length !== 0) {
        type = 'define';
      }
      weeklyData = {
        weekly: {
          ...weekly,
          type: type
        }
      };
    }

    //这一块处理自动生成表达式
    const dayData = {
      day: {
        ...day,
        type: e.target.value
      }
    };
    const data = {
      ...this.state,
      ...dayData,
      ...weeklyData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  onMonthTypeChange = e => {
    const month = this.state.month;

    //这一块处理自动生成表达式
    const monthData = {
      month: {
        ...month,
        type: e.target.value
      }
    };
    const data = {
      ...this.state,
      ...monthData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  onWeeklyTypeChange = e => {
    const weekly = this.state.weekly;
    //这一块处理自动生成表达式
    const weeklyData = {
      weekly: {
        ...weekly,
        type: e.target.value
      }
    };
    //校验
    const day = this.state.day;
    let dayData = { day: { ...day } };
    if (e.target.value !== 'unused') {
      dayData = {
        day: {
          ...day,
          type: 'unused'
        }
      };
    } else {
      let type = 'loop';
      if (day.data.length !== 0) {
        type = 'define';
      }
      dayData = {
        day: {
          ...day,
          type: type
        }
      };
    }

    const data = {
      ...this.state,
      ...weeklyData,
      ...dayData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  onMinDataFromChange = e => {
    const minute = this.state.minute;
    if (e.target.value < 0) {
      e.target.value = 0;
    }
    if (e.target.value > 59) {
      e.target.value = 59;
    }
    const minuteData = {
      minute: {
        ...minute,
        data: {
          ...minute.data,
          from: e.target.value
        }
      }
    };
    const data = {
      ...this.state,
      ...minuteData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  onMinDataToChange = e => {
    const minute = this.state.minute;
    if (e.target.value < 1) {
      e.target.value = 1;
    }
    if (e.target.value > 59) {
      e.target.value = 59;
    }

    const minuteData = {
      minute: {
        ...minute,
        data: {
          ...minute.data,
          to: e.target.value
        }
      }
    };
    const data = {
      ...this.state,
      ...minuteData
    };
    this.setState({
      ...data
    });
    this.getCrontab(data);

  };
  handleMinIndexChange = value => {
    const minute = this.state.minute;
    //这一块处理自动生成表达式
    const minuteData = {
      minute: {
        ...minute,
        data: {
          ...minute.data,
          index: value
        }
      }
    };
    const data = {
      ...this.state,
      ...minuteData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);


  }
  handleHourIndexChange = value => {
    const hour = this.state.hour;

    const hourData = {
      hour: {
        ...hour,
        data: {
          index: value
        }
      }
    };
    const data = {
      ...this.state,
      ...hourData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  handleDayIndexChange = value => {
    const day = this.state.day;

    const dayData = {
      day: {
        ...day,
        data: {
          index: value
        }
      }
    };
    const data = {
      ...this.state,
      ...dayData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  handleMonthIndexChange = value => {
    const month = this.state.month;

    const monthData = {
      month: {
        ...month,
        data: {
          index: value
        }
      }
    };
    const data = {
      ...this.state,
      ...monthData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }
  handleWeeklyIndexChange = value => {
    const weekly = this.state.weekly;

    const weeklyData = {
      weekly: {
        ...weekly,
        data: {
          index: value
        }
      }
    };
    const data = {
      ...this.state,
      ...weeklyData
    };

    this.setState({
      ...data
    });
    this.getCrontab(data);

  }


  minContent = (minute) => {
    const children = [];
    for (let i = 0; i <= 59; i++) {
      children.push(<Option key={i}>{i}</Option>);
    }
    return (
      <div>
        <RadioGroup defaultValue={minute.type} onChange={this.onMinTypeChange}>
          <Radio style={radioStyle} value='loop'>循环&nbsp;&nbsp;
            从第&nbsp;
            <Input
              onChange={this.onMinDataFromChange}
              defaultValue={minute.data.from}
              size='middle'
              min='0'
              max='59'
              type='number'
              style={{ width: 50 }}
              disabled={minute.type !== 'loop'}
            />
            &nbsp;分钟开始，每隔&nbsp;
            <Input
              onChange={this.onMinDataToChange}
              defaultValue={minute.data.to}
              type='number'
              min='1'
              max='59'
              size='middle'
              style={{ width: 50 }}
              disabled={minute.type !== 'loop'}
            />
            &nbsp;分钟执行
          </Radio>
          <Radio style={radioStyle} value='define'>指定&nbsp;&nbsp;
            <Select
              mode='multiple'
              style={{ width: 280 }}
              placeholder=''
              defaultValue={minute.data.index}
              onChange={this.handleMinIndexChange}
              disabled={minute.type !== 'define'}
            >
              {children}
            </Select>
          </Radio>
        </RadioGroup>
      </div>
    );
  }
  hourContent = (hour) => {
    const children = [];
    for (let i = 0; i <= 23; i++) {
      children.push(<Option key={i}>{i}</Option>);
    }
    return (
      <div>
        <RadioGroup defaultValue={hour.type} onChange={this.onHourTypeChange}>
          <Radio style={radioStyle} value='loop'>每小时&nbsp;&nbsp;
          </Radio>
          <Radio style={radioStyle} value='define'>指定&nbsp;&nbsp;
            <Select
              mode='multiple'
              style={{ width: 280 }}
              placeholder=''
              defaultValue={hour.data.index}
              onChange={this.handleHourIndexChange}
              disabled={hour.type !== 'define'}
            >
              {children}
            </Select>
          </Radio>
        </RadioGroup>
      </div>
    );
  }
  dayContent = (day) => {
    const children = [];
    for (let i = 1; i <= 31; i++) {
      children.push(<Option key={i}>{i}</Option>);
    }
    return (
      <div>
        <RadioGroup defaultValue={day.type} value={day.type} onChange={this.onDayTypeChange}>
          <Radio style={radioStyle} value='unused'>不使用天&nbsp;&nbsp;
          </Radio>
          <Radio style={radioStyle} value='loop'>每天&nbsp;&nbsp;
          </Radio>
          <Radio style={radioStyle} value='define' >指定&nbsp;&nbsp;
            <Select
              mode='multiple'
              style={{ width: 280 }}
              placeholder=''
              defaultValue={day.data.index}
              onChange={this.handleDayIndexChange}
              disabled={day.type !== 'define'}
            >
              {children}
            </Select>
          </Radio>
        </RadioGroup>
      </div>
    );
  }
  monthContent = (month) => {
    const children = [];
    for (let i = 1; i <= 12; i++) {
      children.push(<Option key={i}>{i}</Option>);
    }
    return (
      <div>
        <RadioGroup defaultValue={month.type} onChange={this.onMonthTypeChange}>
          <Radio style={radioStyle} value='loop'>每月&nbsp;&nbsp;
          </Radio>
          <Radio style={radioStyle} value='define'>指定&nbsp;&nbsp;
            <Select
              mode='multiple'
              style={{ width: 280 }}
              placeholder=''
              defaultValue={month.data.index}
              onChange={this.handleMonthIndexChange}
              disabled={month.type !== 'define'}
            >
              {children}
            </Select>
          </Radio>
        </RadioGroup>
      </div>
    );
  }
  weeklyContent= (weekly) => {
    const children = [];
    for (let i = 1; i <= 7; i++) {
      children.push(<Option key={i}>星期{i}</Option>);
    }
    return (
      <div>
        <RadioGroup defaultValue={weekly.type} value={weekly.type} onChange={this.onWeeklyTypeChange}>
          <Radio style={radioStyle} value='unused'>不使用星期&nbsp;&nbsp;
          </Radio>
          <Radio style={radioStyle} value='loop'>每星期&nbsp;&nbsp;
          </Radio>
          <Radio style={radioStyle} value='define'>指定&nbsp;&nbsp;
            <Select
              mode='multiple'
              style={{ width: 280 }}
              placeholder=''
              defaultValue={weekly.data.index}
              onChange={this.handleWeeklyIndexChange}
              disabled={weekly.type !== 'define'}
            >
              <Option key='1'>星期一</Option>
              <Option key='2'>星期二</Option>
              <Option key='3'>星期三</Option>
              <Option key='4'>星期四</Option>
              <Option key='5'>星期五</Option>
              <Option key='6'>星期六</Option>
              <Option key='7'>星期日</Option>
            </Select>
          </Radio>
        </RadioGroup>

      </div>
    );
  }

  getCrontab = (data) => {
    const { setFieldsValue } = this.props.form;
    if (data.type === 'define') {
      let crontabArr = [];
      crontabArr[0] = 0;
      //分
      if (data.minute.type === 'loop') {
        crontabArr[1] = `${data.minute.data.from}/${data.minute.data.to}`;
      } else {
        crontabArr[1] = data.minute.data.index.join(',');
      }

      //小时
      if (data.hour.type === 'loop') {
        crontabArr[2] = '*';
      } else {
        crontabArr[2] = data.hour.data.index.join(',');
      }

      //天
      if (data.day.type === 'loop') {
        crontabArr[3] = '*';
      } else if (data.day.type === 'define') {
        crontabArr[3] = data.day.data.index.join(',');
      } else {
        crontabArr[3] = '?';
      }

      //月
      if (data.month.type === 'loop') {
        crontabArr[4] = '*';
      } else {
        crontabArr[4] = data.month.data.index.join(',');
      }

      //周
      if (data.weekly.type === 'loop') {
        crontabArr[5] = '*';
      } else if (data.weekly.type === 'define') {
        crontabArr[5] = data.weekly.data.index.join(',');
      } else {
        crontabArr[5] = '?';
      }

      this.setState({
        crontabExpression: crontabArr.join(' ')
      });

      setFieldsValue({
        crontabUi: { ...data, crontabExpression: crontabArr.join(' ') },
        cron: crontabArr.join(' ')
      });

      this.props.handleClick({ ...data, crontabExpression: crontabArr.join(' ') });
    }
  };
  render() {
    const data = this.props.initialValue || {};
    //const crontabUi = data.crontabUi ? JSON.parse(data.crontabUi) : this.state;
    const crontabUi = this.state;
    const text = <span>Cron表达式的格式：秒 分 时 日 月 周 年(可选)</span>;
    const content = (
      <div>
        <table>
          <tr>
            <th>字段名</th>
            <th>允许的值</th>
            <th>允许的特殊字符</th>
          </tr>
          <tr>
            <td>Seconds (秒)</td>
            <td>0-59</td>
            <td>, - * /</td>
          </tr>
          <tr>
            <td>Minutes (分)</td>
            <td>0-59</td>
            <td>, - * /</td>
          </tr>
          <tr>
            <td>Hours (小时)</td>
            <td>0-23</td>
            <td>, - * /</td>
          </tr>
          <tr>
            <td>Day-of-Month (日)</td>
            <td>1-31</td>
            <td>, - * ? / L W C</td>
          </tr>
          <tr>
            <td>Month (月)</td>
            <td>1-12 or JAN-DEC</td>
            <td>, - * /</td>
          </tr>
          <tr>
            <td>Day-of-Week (周)</td>
            <td>1-7 or SUN-SAT </td>
            <td>, - * ? / L C #</td>
          </tr>
          <tr>
            <td>Year (年,可选字段）</td>
            <td> empty, 1970-2099 </td>
            <td>, - * / </td>
          </tr>
        </table>
        <hr />
        <p> “?”字符：表示不确定的值</p>
        <p> “,”字符：指定数个值</p>
        <p> “-”字符：指定一个值的范围</p>
        <p>“/”字符：指定一个值的增加幅度。n/m表示从n开始，每次增加m</p>
        <p>“L”字符：用在日表示一个月中的最后一天，用在周表示该月最后一个星期X</p>
        <p> “W”字符：指定离给定日期最近的工作日(周一到周五)</p>
        <p>“#”字符：表示该月第几个周X。6#3表示该月第3个周五</p>
      </div>
    );
    const { getFieldDecorator, getFieldValue } = this.props.form;


    const type = getFieldValue('type');
    return (
      <div>
        <FormItem {...tableFormItemLayout} label=''>
          {getFieldDecorator('type', {
            initialValue: crontabUi.type || 'define'
          })(
            <RadioGroup>
              <Radio value='define'>设定</Radio>
              <Radio value='custom'>自定义</Radio>
            </RadioGroup>
          )}
        </FormItem>
        {type === 'define' ? (

          <FormItem {...tableFormItemLayout} label=''>
            {getFieldDecorator('crontabUi', {
            })(
              <div className={styles.cardContainer}>
                <Tabs type='card'>
                  <TabPane tab='分钟' key='minute'>
                    {this.minContent(crontabUi.minute)}
                  </TabPane>
                  <TabPane tab='小时' key='hour'>
                    {this.hourContent(crontabUi.hour)}
                  </TabPane>
                  <TabPane tab='天' key='day'>
                    {this.dayContent(crontabUi.day)}
                  </TabPane>
                  <TabPane tab='月' key='month'>
                    {this.monthContent(crontabUi.month)}
                  </TabPane>
                  <TabPane tab='星期' key='weekly'>
                    {this.weeklyContent(crontabUi.weekly)}
                  </TabPane>
                </Tabs>
                <div style={{ marginTop: 10 }}>
                  <Input disabled={true} style={{ float: 'left', width: 400 }} defaultValue={data.crontab} value={this.state.crontabExpression} />
                </div>
              </div>
            )}
          </FormItem>
        ) : (
          <FormItem
            {...tableFormItemLayout}
            label=''
            extra='例&quot;0 0 12 ? * WED&quot;表示在每星期三下午12:00执行'
          >
            {getFieldDecorator('cron', {})(
              <div style={{ position: 'relative' }}>
                <Input placeholder='请输入 crontab 表达式' defaultValue={data.crontab} />
                <Popover
                  placement='right'
                  title={text}
                  content={content}
                  trigger='click'
                >
                  <Icon
                    type='info-circle-o'
                    style={{
                      fontSize: 16,
                      color: '#f6ae68',
                      position: 'absolute',
                      right: -22,
                      top: 10
                    }}
                  />
                </Popover>
              </div>
            )}
          </FormItem>
        )}
      </div>
    );
  }
}
