import React from 'react';
import { Input, DatePicker, Calendar, Icon, Radio, Row, Col, Checkbox, Button, Select } from 'antd';
const CheckboxGroup = Checkbox.Group;
const RadioGroup = Radio.Group;
const { RangePicker } = DatePicker;
import SearchTree from './tree';
import styles from './search.less';
import Item from 'antd/lib/list/Item';

/**
 * 数据模型
 * 搜索项searchKeys[]，{label:'',key:''}
 * 搜索值searchValues{}，{key:{}}
 * 初始值/快速搜索：[{
        key: 'usage',
        keyLabel: '用途',
        value: 'production',
        valueLabel: '生产'        
      }]
 */

export default class Container extends React.Component {
  constructor(props) {
    super(props);
    this.myRef = React.createRef();
    this.state = {
      historyList: [],
      showHistoryPanel: false,
      KeyDisplay: 'none',
      valueDisplay: 'none',
      left: 0,
      checkboxValue: [],
      searchKey: '',
      inputValue: '',
      section1: '',
      searchInputValue: '',
      searchAreaValue: '',
      searchList: [
        // { key: 'radio', KeyLabel: '单选', valueLabel:'A1', value: 'a1' }
      ],
      searchKeys: props.searchKeys,
      searchValues: props.searchValues,
      cascader_children: [],
      cascader_label: [],
      cascader_value: [],
      cascader_label_aysnc: [],
      cascader_value_aysnc: []
    };
  }

 
  componentDidMount() {
    document.body.addEventListener('click', this.handleDocumentClick);
    //设置初始值
    const { initialSearchList } = this.props;
    if (initialSearchList) {
      this.onSubmitFromHistory(initialSearchList);
    }
  }

  componentWillUnmount() {
    document.body.removeEventListener('click', this.handleDocumentClick);
  }

  hidePanel = () => {
    this.setState({
      KeyDisplay: 'none',
      valueDisplay: 'none',
      showHistoryPanel: false
    });
  };


  //提交数据
  submit = (e, type) => {
    this.hidePanel();
    if (type === 'input' && e.target.value) {
      return this.props.onSearch({
        keyword: e.target.value
      });
    }
    const historyList = this.state.historyList;
    const data = {};
    const searchList_copy = JSON.parse(JSON.stringify(this.state.searchList));
    historyList.push(searchList_copy);
    this.setState({
      historyList
    });
    searchList_copy.forEach(item => {
      data[item.key] = item.value;
    });
    this.props.onSearch(data);
  };

  //快速搜索策略，添加搜索项，已存在则替换
  onSubmitFromQuikSearch = (data) => {  
    let { searchList } = this.state;  
    const list = JSON.parse(JSON.stringify(searchList));
    list.push(data);
    for (let i = 0; i < searchList.length; i ++) {
      if (searchList[i].key === data.key) {
        list.splice(i, 1);
      }
    }
    this.onSubmitFromHistory(list);
  }

  onSubmitFromHistory = (list) => {
    this.hidePanel();
    const data = {};
    let { searchKeys, searchList } = this.state;
    const copy_searchList = JSON.parse(JSON.stringify(searchList));
    const copy_searchKeys = JSON.parse(JSON.stringify(searchKeys));

    //并集, 当前的searchList + 当前的searchKeys = 总的searchKeys
    const copy_searchList_to_searchKeys = copy_searchList.map(item => {
      return {
        label: item.keyLabel, key: item.key
      };
    });
    const union = [ ...copy_searchKeys, ...copy_searchList_to_searchKeys ];
    //差集  总的searchKeys - list =  searchKeys
    const difference = union.filter(item => !list.some(i => i.key === item.key));
    this.setState({
      searchList: list,
      searchKeys: difference
    });
    list.forEach(item => {
      data[item.key] = item.value;
    });
    this.props.onSearch(data);
  };

  //点击输入框
  dropDown = (e) => {
    //KeyMenu位置重置
    this.setState({
      left: 0
    });
    const { searchList } = this.state;
    if (searchList.length === 0) {
      this.setState({
        valueDisplay: 'none',
        KeyDisplay: 'block'
      });
    } else {
      if (!searchList[searchList.length - 1].value) {
        //没有选择值时，显示值选框
        this.setState({
          KeyDisplay: 'none',
          valueDisplay: 'block'
        });
      } else {
        this.setState({
          valueDisplay: 'none',
          KeyDisplay: 'block'
        });
      }
    }
  };

  //搜索框中输入搜索Key
  searchKey = (e) => {
    this.setState({
      inputValue: e.target.value
    });
  };


  //删除某一项
  remove = (e, index, item) => {
    let { searchList, searchKeys } = this.state;
    searchList.splice(index, 1); //显示列表去掉一项
    searchKeys.splice(index, 0, { label: item.keyLabel, key: item.key }); //搜索项增加一项
    this.setState({
      searchKey: '',
      searchList,
      searchKeys,
      cascader_children: []
    });
    this.myRef.current.focus();
  };

  //清空全部
  clear = () => {
    let { searchList, searchKeys } = this.state;
    searchList.forEach(item => {
      searchKeys.push({ label: item.keyLabel, key: item.key });
    });
    this.setState({
      searchKeys: searchKeys,
      searchList: [],
      searchKey: '',
      inputValue: '',
      cascader_children: []
    });
    this.myRef.current.focus();
  };

  //点击搜索项
  onSearchKeyChange = (e) => {
    const value = e.target.value;
    let { searchKeys, searchList } = this.state;
    searchKeys.forEach((item, index) => {
      if (item.key === value) {
        searchKeys.splice(index, 1);
        searchList.push({ keyLabel: item.label, key: item.key });
      }
    });
    this.setState({
      searchKey: value,
      KeyDisplay: 'none',
      valueDisplay: 'block',
      showHistoryPanel: false,
      inputValue: '',
      searchInputValue: '',
      searchAreaValue: '',
      searchList,
      searchKeys
    });
  };

  //复选框点击确定按钮触发，先存值
  checkboxChange = (value) => {
    this.setState({
      checkboxValue: value
    });
  };

  //Section
  saveSectionValue = (e) => {
    this.setState({
      section1: e.target.value
    });
  };

  //点击搜索值
  onSearchValueChange = (e, type, options) => {
    let value;
    let label;
    switch (type) {
    case 'radio':
      const v = JSON.parse(e.target.value);
      value = v.value;
      label = v.label;
      
      break;
    case 'checkbox':
      const v2 = this.state.checkboxValue.map(data => JSON.parse(data));
      value = v2.map(data => data.value);
      label = v2.map(data => data.label).join(',');
      break;
    case 'input':
      value = label = e.target.value;
      break;
    case 'textArea':
      value = label = this.state.searchAreaValue;
      break;
    case 'number':
      value = label = e.target.value;
      break;
    case 'rangePicker':
      const t0 = e[0].format(options.format || 'YYYY-MM-DD HH:mm:ss');
      const t1 = e[1].format(options.format || 'YYYY-MM-DD HH:mm:ss');
      value = [ t0, t1 ];       
      label = t0 + ' ~ ' + t1;
      break;
    case 'datePicker':
      value = label = e.format(options.format || 'YYYY-MM-DD HH:mm:ss');    
      break;
    case 'calendar':
      value = label = e.format(options.format || 'YYYY-MM-DD');    
      break;
    case 'tree':
      value = e.id;
      label = e.name;
      break;
    case 'cascader_ysnc':
      const v3 = JSON.parse(e.target.value);
      value = [ this.state.cascader_value, v3.id ];
      label = this.state.cascader_label + ' / ' + v3.name;
      break;
    case 'cascader_aysnc':
      const v4 = JSON.parse(e.target.value);
      value = [ ...this.state.cascader_value_aysnc, v4.value ];
      label = [ ...this.state.cascader_label_aysnc, v4.label ].join(' / ');
      break;
    case 'section':
      const v5 = e.target.value;
      value = [ this.state.section1, v5 ];
      label = [ this.state.section1, v5 ].join(' - ');
      break;
    default:
      value = '';
      label = '';
    }

    let { searchKey, searchList } = this.state;
    searchList.forEach(item => {
      if (item.key === searchKey) {
        item.value = value;
        item.valueLabel = label;
      }
    });

    this.setState({
      searchList,
      valueDisplay: 'none',
      KeyDisplay: 'block' //用户操作方便性来讲，选完一个值后显示key框比较方便
    });
    this.myRef.current.focus();
  };

  //修改值，目前只有输入框的值才有回显
  modifyValue = (item) => {
    let width = this.ulContainer.clientWidth - this.inputContainer.clientWidth - this[`inputLi${item.key}`].offsetLeft;
    this.setState({
      valueDisplay: 'block',
      KeyDisplay: 'none',
      searchKey: item.key,
      searchInputValue: item.value,
      searchAreaValue: Item.value,
      left: -width //重定位
    });
  }

  //打开历史记录面板
  openHistoryPanel = () => {
    this.setState({
      showHistoryPanel: !this.state.showHistoryPanel,
      valueDisplay: 'none',
      KeyDisplay: 'none'
    });
  };
  //关闭历史记录面板
  closeHistoryPanel = () => {
    this.setState({
      showHistoryPanel: false
    });
  };

  //清除历史记录
  clearHistory = () => {
    this.setState({
      historyList: []
    });
  };

  onSearchInputChange = (e) => {
    //输入框的值需要可控制清除
    this.setState({
      searchInputValue: e.target.value
    });
  };
  onSearchAreaChange = (e) => {
    this.setState({
      searchAreaValue: e.target.value
    });
  };

  showChildren = (e) => {
    const v = JSON.parse(e.target.value);
    let { searchKey, searchList } = this.state;
    searchList.forEach(item => {
      if (item.key === searchKey) {
        item.value = [v.value]; //值为数组形式
        item.valueLabel = v.label;
      }
    });
    //保存值，选二级时需要用到
    this.setState({
      cascader_children: v.children,
      cascader_label: v.label,
      cascader_value: v.value
    });
    //只想选一级
    this.myRef.current.focus();
  };

  showSecond = (e, cb) => {
    const v = JSON.parse(e.target.value);
    let { searchKey, searchList } = this.state;
    searchList.forEach(item => {
      if (item.key === searchKey) {
        item.value = [v.value]; //值为数组形式
        item.valueLabel = [v.label];
      }
    });
    cb(v.value);
    this.setState({
      cascader_label_aysnc: [v.label],
      cascader_value_aysnc: [v.value]
    });
    this.myRef.current.focus();
  };

  showThird = (e, grade, type, cb) => {
    //只有两级
    if (grade === 2) {
      this.onSearchValueChange(e, type);
    } else {
      //第三级
      const v = JSON.parse(e.target.value);
      let { searchKey, searchList } = this.state;
      searchList.forEach(item => {
        if (item.key === searchKey) {
          item.value[1] = v.value; //值为数组形式
          item.valueLabel[1] = v.label;
          this.setState({
            cascader_label_aysnc: item.valueLabel,
            cascader_value_aysnc: item.value
          });
        }
      });
      cb(v.value);
      this.myRef.current.focus();
    }
  };
  showFourth = (e, grade, type, cb) => {
    //只有3级
    if (grade === 3) {
      this.onSearchValueChange(e, type);
    } else {
      //第四级
      const v = JSON.parse(e.target.value);
      let { searchKey, searchList } = this.state;
      searchList.forEach(item => {
        if (item.key === searchKey) {
          item.value[2] = v.value; //值为数组形式
          item.valueLabel[2] = v.label;
          this.setState({
            cascader_label_aysnc: item.valueLabel,
            cascader_value_aysnc: item.value
          });
        }
      });
      cb(v.value);
      this.myRef.current.focus();
    }
  };

  
  handleDocumentClick = (e) => {
    if (this.dropDownKeyMenu && this.dropDownKeyMenu.contains(e.target)) {
      return;
    }
    if (this.dropDownValueMenu && this.dropDownValueMenu.contains(e.target)) {
      return;
    }
    if (this.dropDownHistoryMenu && this.dropDownHistoryMenu.contains(e.target)) {
      return;
    }
    //antd calendar 绝对定位在dom外面
    const calendarEle = document.getElementsByClassName('ant-calendar-picker-container')[0];
    if (calendarEle && calendarEle.contains(e.target)) {
      return;
    }
    this.hidePanel();
  }

  render() {
    const searchValues = this.state.searchValues[this.state.searchKey] || {};
    return (
      <div id='idcosSearch' className={styles['idcos-search-container']}>
        <div className={styles['filtered-search-box']}>
          <div className='filtered-search-history-dropdown-wrapper' ref={node => this.dropDownHistoryMenu = node} >
            <button onClick={this.openHistoryPanel} className='filtered-search-history-dropdown-toggle-button' type='button'>
              <Icon type='clock-circle' theme='outlined' />
            </button>
            {
              this.state.showHistoryPanel &&
                <div className='dropdown-menu filtered-search-history-dropdown' >
                  <div className='dropdown-title'>
                    <span>最近搜索项</span>
                    <Icon onClick={this.closeHistoryPanel} className='dropdown-menu-close' type='close' theme='outlined' />
                  </div>
                  <div className='dropdown-content'>
                    <div>
                      <ul className='filtered-search-history-dropdown-content'>
                        {
                          this.state.historyList.length > 0 ? this.state.historyList.map(arr =>
                            <li onClick={() => this.onSubmitFromHistory(arr)}>
                              {
                                arr.length > 0 ? arr.map(item =>
                                  <span className='filtered-search-history-dropdown-token'>
                                    <span className='name'>{item.keyLabel}：</span>
                                    <span className='value'>{item.valueLabel}</span>
                                  </span>) : <span>无搜索条件</span>
                              }
                            </li>
                          ) : <li className='history-none'>您没有最近搜索项</li>
                        }
                        {
                          this.state.historyList.length > 0 &&
                          <li className='history-clear' onClick={this.clearHistory}>清除搜索记录</li>
                        }
                      </ul>
                    </div>
                  </div>
                </div>
            }
          </div>
          <div className={styles['filtered-search-box-input-container']}>
            <div className={styles['scroll-container']}>
              <ul className={styles['tokens-container']} ref={node => this.ulContainer = node}>
                {
                  this.state.searchList.map((item, index) => <li key={item.key} ref={node => this[`inputLi${item.key}`] = node} className={styles['filtered-search-token']}>
                    <div className={styles['selectable']}>
                      <div onClick={() => this.modifyValue(item)} >
                        {item.keyLabel}：
                          
                        {item.valueLabel}
                      </div>
                      <div className={styles['remove-token']} onClick={(e) => this.remove(e, index, item)}>
                        <Icon type='close' />
                      </div>
                    </div>
                  </li>)
                }
                <li className={styles['input-token']} ref={node => this.inputContainer = node}>
                  <Input ref={this.myRef} onChange={this.props.onChange} enterButton={true} onPressEnter={(e) => this.submit(e, 'input')} onClick={this.dropDown} autoComplete='off' autoFocus={true} className={styles['form-control']} placeholder={this.props.placeholder || '请选择...'} />
                  <div ref={node => this.dropDownKeyMenu = node} style={{ display: this.state.KeyDisplay }} className='dropdown-menu dropdown-key-menu'>
                    <RadioGroup style={{ width: '100%' }} value={this.state.searchKey} onChange={this.onSearchKeyChange}>
                      <Row>
                        <Col className='search-btn' onClick={this.submit}>
                          <Icon type='search' style={{ marginRight: 8 }} />
                                        按回车键或点击搜索
                        </Col>
                        {
                          this.state.searchKeys.length > 0 &&
                                        (this.state.inputValue ? this.state.searchKeys.filter(data => data.label.indexOf(this.state.inputValue) !== -1) : this.state.searchKeys).map(item => <Col span={24}><Radio value={item.key}>{item.label}</Radio></Col>)
                        }
                      </Row>
                    </RadioGroup>
                  </div>
                  <div ref={node => this.dropDownValueMenu = node} style={{ left: this.state.left, display: this.state.valueDisplay }} className='dropdown-menu'>
                    {
                      searchValues.type === 'cascader_aysnc' &&
                  <Row>
                    <Col span={24 / searchValues.grade} className='idcos_search_cascader'>
                      <RadioGroup style={{ width: '100%' }} onChange={(e) => this.showSecond(e, searchValues.getSecondList)}>
                        <Row>
                          {searchValues.list.map(item => <Col title={item.label} className={styles['col']} span={24}><Radio value={JSON.stringify(item)}>{item.label}</Radio></Col>)}
                        </Row>
                      </RadioGroup>
                    </Col>
                    <Col span={24 / searchValues.grade}>
                      <RadioGroup style={{ width: '100%' }} onChange={(e) => this.showThird(e, searchValues.grade, searchValues.type, searchValues.getThirdList)}>
                        <Row>
                          {(searchValues.secondList || []).map(item => <Col className={styles['col']} span={24}><Radio value={JSON.stringify(item)}>{item.label}</Radio></Col>)}
                        </Row>
                      </RadioGroup>
                    </Col>
                    {
                      searchValues.thirdList &&
                      <Col span={24 / searchValues.grade}>
                        <RadioGroup style={{ width: '100%' }} onChange={(e) => this.showFourth(e, searchValues.grade, searchValues.type, searchValues.getFourthList)}>
                          <Row>
                            {(searchValues.thirdList || []).map(item => <Col className={styles['col']} span={24}><Radio value={JSON.stringify(item)}>{item.label}</Radio></Col>)}
                          </Row>
                        </RadioGroup>
                      </Col>
                    }
                    {
                      searchValues.fourthList &&
                      <Col span={24 / searchValues.grade}>
                        <RadioGroup style={{ width: '100%' }} onChange={(e) => this.onSearchValueChange(e, searchValues.type)}>
                          <Row>
                            {(searchValues.fourthList || []).map(item => <Col className={styles['col']} span={24}><Radio value={JSON.stringify(item)}>{item.label}</Radio></Col>)}
                          </Row>
                        </RadioGroup>
                      </Col>
                    }
                  </Row>
                    }
                    {
                      searchValues.type === 'radio' &&
                  <RadioGroup style={{ width: '100%' }} onChange={(e) => this.onSearchValueChange(e, searchValues.type)}>
                    <Row>
                      {searchValues.list.map(item => <Col className={styles['col']} span={24}><Radio value={item.label ? JSON.stringify(item) : JSON.stringify({ label: item, value: item })}>{item.label ? item.label : item}</Radio></Col>)}
                    </Row>
                  </RadioGroup>
                    }
                    {
                      searchValues.type === 'checkbox' &&
                  <div>
                    <CheckboxGroup key={this.state.searchKey} style={{ width: '100%' }} onChange={this.checkboxChange}>
                      <Row>
                        {searchValues.list.map(item => <Col className={styles['col']} span={24}><Checkbox value={item.label ? JSON.stringify(item) : JSON.stringify({ label: item, value: item })}>{item.label ? item.label : item}</Checkbox></Col>)}
                      </Row>
                    </CheckboxGroup>
                    <div className='checkbox-confirm' onClick={(e) => this.onSearchValueChange(e, searchValues.type)}>点击确认</div>
                  </div>
                    }
                    {
                      searchValues.type === 'input' &&
                  <Row>
                    <Col>
                      <Input autoFocus={true} value={this.state.searchInputValue} onChange={this.onSearchInputChange} placeholder={searchValues.placeholder || '请输入，按回车键确定'} onPressEnter={(e) => this.onSearchValueChange(e, searchValues.type)} />
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'textArea' &&
                  <Row>
                    <Col>
                      <Input.TextArea rows={3} autoFocus={true} placeholder={searchValues.placeholder} onChange={this.onSearchAreaChange} value={this.state.searchAreaValue} />
                      <Button size='small' style={{ float: 'right', marginTop: 8 }} type='primary' onClick={(e) => this.onSearchValueChange(e, searchValues.type)}>确认</Button>
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'rangePicker' &&
                  <Row>
                    <Col>
                      <RangePicker
                        showTime={searchValues.showTime !== 'true'}
                        style={{ width: '100%' }}
                        format={searchValues.format || 'YYYY-MM-DD HH:mm:ss'}
                        placeholder={searchValues.placeholder || [ '开始时间', '结束时间' ]}
                        onOk={(e) => this.onSearchValueChange(e, searchValues.type, searchValues)}
                      />
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'datePicker' &&
                  <Row>
                    <Col>
                      <DatePicker
                        showToday={false}
                        showTime={searchValues.showTime !== 'true'}
                        style={{ width: '100%' }}
                        format={searchValues.format || 'YYYY-MM-DD HH:mm:ss'}
                        placeholder={searchValues.placeholder || ['选择时间']}
                        onOk={(e) => this.onSearchValueChange(e, searchValues.type, searchValues)}
                      />
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'calendar' &&
                  <Row>
                    <Col>

                      <Calendar fullscreen={false} onSelect={(e) => this.onSearchValueChange(e, searchValues.type, searchValues)} />
                    
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'number' &&
                  <Row>
                    <Col>
                      <Input type='number' autoFocus={true} placeholder={searchValues.placeholder || '请输入，按回车键确定'} onPressEnter={(e) => this.onSearchValueChange(e, searchValues.type)} />
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'tree' &&
                  <Row>
                    <Col>
                      <SearchTree
                        onSelect={(e) => this.onSearchValueChange(e, searchValues.type)}
                        treeData={searchValues.data}
                      />
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'cascader_ysnc' &&
                  <Row>
                    <Col span='12' className='idcos_search_cascader'>
                      <RadioGroup style={{ width: '100%' }} onChange={(e) => this.showChildren(e)}>
                        <Row>
                          {searchValues.list.map(item => <Col title={item.label} className={styles['col']} span={24}><Radio value={JSON.stringify(item)}>{item.label}</Radio></Col>)}
                        </Row>
                      </RadioGroup>
                    </Col>
                    <Col span='12'>
                      <RadioGroup style={{ width: '100%' }} onChange={(e) => this.onSearchValueChange(e, searchValues.type)}>
                        <Row>
                          {(this.state.cascader_children || []).map(item => <Col className={styles['col']} span={24}><Radio value={JSON.stringify(item)}>{item.name}</Radio></Col>)}
                        </Row>
                      </RadioGroup>
                    </Col>
                  </Row>
                    }
                    {
                      searchValues.type === 'section' &&
                  <Row>
                    <Col span='10'>
                      <Input autoFocus={true} type='number' max={searchValues.max} min={searchValues.min || 0} step={searchValues.step || 0.01} onBlur={this.saveSectionValue} />
                    </Col>
                    <Col span='4' className='section-middle-line'>
                                        -
                    </Col>
                    <Col span='10'>
                      <Input type='number' max={searchValues.max} min={searchValues.min || 0} step={searchValues.step || 0.01} onPressEnter={(e) => this.onSearchValueChange(e, searchValues.type)} />
                    </Col>
                  </Row>
                    }
                  </div>
           
                </li>

              </ul>
            </div>
          </div>

          {
            this.props.quikSearch && 
            <div className={styles['quik-search']}>
              <div className={styles['quik-search-tags']}>
                {
                  this.props.quikSearch.map(q => <div onClick={() => this.onSubmitFromQuikSearch(q)}>{q.valueLabel}</div>)
                }
              </div>
              {/* <Select className={styles['quik-search-select']} mode='multiple' onChange={this.onSubmitFromHistory} defaultValue={[ 1, 12 ]}>
                <Option key={1} value={1}>{1}</Option>
                <Option key={12} value={12}>{12}</Option>
                <Option key={13} value={13}>{14}</Option>
              </Select>            */}
            </div> 
          }     
          {
            this.state.searchList.length > 0 &&
              <button className={styles['clear-search']} type='button' onClick={this.clear}>
                <Icon type='close' />
              </button>
          }
        </div>
      </div>
    );
  }
}

