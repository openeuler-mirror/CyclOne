import React from 'react';
export default class Normal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      myCodeMirror: null
    };
    this.myRef = React.createRef();
  }

  componentDidMount() {
    const myCodeMirror = CodeMirror(this.myRef.current, {
      value: this.props.data,
      mode: 'shell',
      theme: 'erlang-dark',
      readOnly: true,
      styleActiveLine: true,
      lineNumbers: true,
      highlightSelectionMatches: true
    });
    this.setState({
      myCodeMirror
    });
  }
  UNSAFE_componentWillReceiveProps(nextPorps) {
    const myCodeMirror = this.state.myCodeMirror;
    if (myCodeMirror) {
      myCodeMirror.setValue(nextPorps.data);
    }
  }
  componentWillUnmount() {
    this.setState({
      myCodeMirror: null
    });
  }

  render() {
    return (
      <div className=''>
        <div ref={this.myRef} />
      </div>
    );
  }
}
