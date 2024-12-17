import * as React from 'react';
import { findDOMNode } from 'react-dom';
import {
  DragSource,
  DropTarget,
  ConnectDropTarget,
  ConnectDragSource,
  DropTargetMonitor,
  DropTargetConnector,
  DragSourceConnector,
  DragSourceMonitor
} from 'react-dnd';
import { XYCoord } from 'dnd-core';
import Template from './template';


const cardSource = {
  beginDrag(props) {
    return {
      data: props.data,
      index: props.index
    };
  }
};

const cardTarget = {
  hover(props, monitor, component) {
    if (!component) {
      return null;
    }
    const dragIndex = monitor.getItem().index;
    const hoverIndex = props.index;
    const dragData = props.data;

    // Don't replace items with themselves
    if (dragIndex === hoverIndex) {
      return;
    }

    // Determine rectangle on screen
    const hoverBoundingRect = (findDOMNode(component)).getBoundingClientRect();

    // Get vertical middle
    const hoverMiddleY = (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2;

    // Determine mouse position
    const clientOffset = monitor.getClientOffset();

    // Get pixels to the top
    const hoverClientY = (clientOffset).y - hoverBoundingRect.top;

    // Only perform the move when the mouse has crossed half of the items height
    // When dragging downwards, only move when the cursor is below 50%
    // When dragging upwards, only move when the cursor is above 50%
    // Dragging downwards
    if (dragIndex < hoverIndex && hoverClientY < hoverMiddleY) {
      return;
    }

    // Dragging upwards
    if (dragIndex > hoverIndex && hoverClientY > hoverMiddleY) {
      return;
    }

    // Time to actually perform the action
    props.moveCard(dragIndex, hoverIndex, dragData);

    // Note: we're mutating the monitor item here!
    // Generally it's better to avoid mutations,
    // but it's good here for the sake of performance
    // to avoid expensive index searches.
    monitor.getItem().index = hoverIndex;
  }
};



@DropTarget('card', cardTarget, (connect) => ({
  connectDropTarget: connect.dropTarget()
}))
@DragSource(
  'card',
  cardSource,
  (connect, monitor) => ({
    connectDragSource: connect.dragSource(),
    isDragging: monitor.isDragging()
  }),
)

class Card extends React.Component {
  render() {
    const {
      isDragging,
      connectDragSource,
      connectDropTarget,
      disabled
    } = this.props;

    const opacity = isDragging ? 0.5 : 1;

    if (disabled) {
      return (
        <div style={{ opacity }} className='configure-card'>
          <span className='configure-index'>{this.props.index + 1}</span>
          <div className='configure-content'>
            <Template {...this.props} />
          </div>
          <div className='clearfix' />
        </div>
      );
    }
    return (
      connectDragSource &&
      connectDropTarget &&
      <div style={{ opacity }} className='configure-card'>
        {connectDragSource(
          connectDropTarget(
            <span className='configure-index' title='拖动图标可进行排序操作'>{this.props.index + 1}</span>
          ),
        )}
        <div className='configure-content'>
          <Template {...this.props} />
        </div>
        <div className='clearfix' />
      </div>

    );
  }
}

export default Card;
