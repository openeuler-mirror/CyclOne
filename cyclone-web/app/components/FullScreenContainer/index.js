/**
 * 方便的将 Component 丢到全屏范围
 *
 * Created by zhangrong on 17/4/25.
 */

import React, { Component } from 'react';
import ReactDOM from 'react-dom';

import uuid from 'common/uuid';

class FullScreenContainer extends Component {
  constructor(props) {
    super(props);
    this.state = {
      divId: uuid()
    };
  }
  componentDidMount() {
    this.renderComponent(this);
  }
  componentDidUpdate() {
    this.renderComponent(this);
  }
  componentWillUnmount() {
    const _container = this._container;
    ReactDOM.unmountComponentAtNode(_container);
    document.body.removeChild(_container);
  }
  renderComponent = (instance, componentArg, ready) => {
    if (!instance._container) {
      instance._container = this.getContainer(instance);
    }
    let component;
    if (instance.getComponent) {
      component = instance.getComponent(componentArg);
    } else {
      component = this.getComponent(instance, componentArg);
    }
    ReactDOM.unstable_renderSubtreeIntoContainer(
      instance,
      component,
      instance._container,
      function callback() {
        instance._component = this;
        if (ready) {
          ready.call(this);
        }
      }
    );
  };

  getContainer = instance => {
    if (instance.props.getContainer) {
      return instance.props.getContainer();
    }

    const container = document.createElement('div');
    container.setAttribute('id', this.state.divId);
    container.style.position = 'absolute';
    container.style.zIndex = 999;
    container.style.top = '0';
    container.style.bottom = '0';
    container.style.left = '0';
    container.style.right = '0';
    container.style.width = '100%';
    container.style.height = '100%';

    document.body.appendChild(container);
    return container;
  };

  getComponent = (instance, extra) => {
    return <div>{React.Children.only(this.props.children)}</div>;
  };

  render() {
    return null;
  }
}

export default FullScreenContainer;
