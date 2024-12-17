import * as React from 'react';
import { DragDropContext } from 'react-dnd';
import HTML5Backend from 'react-dnd-html5-backend';
import Card from './card';
const update = require('immutability-helper');
import { Popover, Button, Radio } from 'antd';
const RadioGroup = Radio.Group;
import { category, getTempData, getFormData } from './data';

@DragDropContext(HTML5Backend)
export default class Container extends React.Component {
  constructor(props) {
    super(props);
    this.moveCard = this.moveCard.bind(this);
    this.state = {
      value: '',
      cards: this.getCardsDone().cards || []
    };
    this.uuid = this.getCardsDone().uuid || 0;
  }

  //cards赋值
  getCardsDone = () => {
    const cards = this.props.cards || [];
    const copyCards = [];
    if (cards.length > 0) {
      cards.map((card, index) => {
        copyCards[index] = getTempData()[card.category];
        if (card.category === 'raid') {
          if (card.metadata.controller_index === '0') {
            copyCards[index] = getTempData()['raid0'];
          } else if (card.metadata.controller_index === '1') {
            copyCards[index] = getTempData()['raid1'];
          } else if (card.metadata.controller_index === '2') {
            copyCards[index] = getTempData()['raid2'];
          } else if (card.metadata.controller_index === '3') {
            copyCards[index] = getTempData()['raid3'];
          }
        }
        copyCards[index].initialValue = card;
        copyCards[index].uuid = index;
      });
      this.props.getOrder(copyCards);
    }
    return {
      cards: copyCards,
      uuid: cards.length
    };
  };

  onCategoryChange = (e) => {
    const value = e.target.value;
    const { cards } = this.state;
    const pushData = getTempData()[value];
    this.uuid++;
    pushData.uuid = this.uuid;
    cards.push(pushData);
    this.setState({
      value,
      cards
    });
    this.props.getOrder(cards);
    // this.setState({ value: '' });
  };


  deleteConfig = (index) => {
    const { cards } = this.state;
    cards.splice(index, 1);
    this.setState({ cards });
    this.props.getOrder(cards);
  };

  copyConfig = (data) => {
    // 复制的时候给数据的initialValue赋值表单的值
    const initialValue = getFormData([data], this.props.form.getFieldsValue());
    data.initialValue = initialValue[0];
    const { cards } = this.state;
    const copyData = JSON.parse(JSON.stringify(data));
    this.uuid++;
    copyData.uuid = this.uuid;
    cards.push(copyData);
    this.setState({ cards });
    this.props.getOrder(cards);
  };

  handleClick = () => {
    this.setState({ value: '' });
  };
  render = () => {
    const { cards } = this.state;
    const radioStyle = {
      display: 'block',
      height: '30px',
      lineHeight: '30px'
    };

    const content = (
      <RadioGroup onChange={this.onCategoryChange} value={this.state.value}>
        {
          category.map(data => {
            return (
              <Radio style={radioStyle} value={data.value}>
                {data.category}
              </Radio>
            );
          })
        }
      </RadioGroup>
    );

    return (
      <div>
        {cards.map((card, i) => (
          <Card
            key={i}
            index={i}
            data={card}
            moveCard={this.moveCard}
            deleteConfig={this.deleteConfig}
            copyConfig={this.copyConfig}
            form={this.props.form}
            dictionaries={this.props.dictionaries}
            firmwares={this.props.firmwares}
            disabled={this.props.disabled}
          />
        ))}
        {
          !this.props.disabled &&
          <div className='addConfigBtn'>
            <Popover placement='bottom' content={content} trigger='click'>
              <Button icon='plus' type='primary' onClick={this.handleClick}>添加配置项</Button>
            </Popover>
          </div>
        }
      </div>
    );
  };

  moveCard = (dragIndex, hoverIndex) => {
    const { cards } = this.state;
    const dragCard = cards[dragIndex];
    this.setState(
      update(this.state, {
        cards: {
          $splice: [[ dragIndex, 1 ], [ hoverIndex, 0, dragCard ]]
        }
      }),
    );
    this.props.getOrder(this.state.cards);
  }
}
