import React from 'react';
import Patrol from './patrol';
import Normal from './normal';
import Popup from '../popup';
const TYPES = {
  normal: Normal,
  patrol: Patrol
};

export class ExecuteResult extends React.Component {
  constructor(props) {
    super(props);
    this.threadHandle = null;
    this.state = {
      logInfo: ''
    };
  }
  componentDidMount() {
    if (this.props.getData) {
      this.initThread();
    }
  }

  initThread = () => {
    const { record } = this.props;
    const getLog = () => {
      this.props.getData(record, this);
    };
    setTimeout(getLog, 100);  //componentDidMount函数给编辑器立即复制，样式会乱
    this.threadHandle = setInterval(getLog, 5000);
  };

  componentWillUnmount() {
    clearInterval(this.threadHandle);
  }

  render() {
    const { type, data } = this.props;
    const Co = TYPES[type] || TYPES.normal;
    return (
      <div className='execute-result'>
        <Co data={data || this.state.logInfo} />
      </div>
    );
  }
}
export default function open(props) {
  Popup.open({
    title: `${props.title || '执行结果详情'} `,
    maskClosable: false,
    width: 800,
    onCancel: () => {
      Popup.close();
      if (props.getData) {
        props.reload();
      }
    },
    content: <ExecuteResult {...props} />
  });
}
