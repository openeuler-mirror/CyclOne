import React from 'react';
import SearchForm from '../index';
function getWithArgs(a, b) {
  console.log();
}
const mockData = {
  "id": '0',
  "name": "全部",
  "capacity": 99999,
  "children": [
    {
      "id": '1',
      "name": "中国",
      "capacity": 1000,
      "used": 0,
      "free": 1000,
      "children": [
        {
          "id": '3',
          "name": "浙江",
          "capacity": 500,
          "used": 0,
          "free": 500,
          "children": [
            {
              "id": '4',
              "name": "杭州",
              "capacity": 100,
              "used": 0,
              "free": 100,
              "children": [ ]
            }
          ]
        }
      ]
    },
    {
      "id": '2',
      "name": "美国",
      "capacity": 1000,
      "used": 0,
      "free": 1000,
      "children": [
        {
          "id": '51',
          "name": "杭州",
          "capacity": 100,
          "used": 0,
          "free": 100,
          "children": [
            {
              "id": '11',
              "name": "中国",
              "capacity": 1000,
              "used": 0,
              "free": 1000,
              "children": [
                {
                  "id": '31',
                  "name": "浙江",
                  "capacity": 500,
                  "used": 0,
                  "free": 500,
                  "children": [
                    {
                      "id": '41',
                      "name": "杭州",
                      "capacity": 100,
                      "used": 0,
                      "free": 100,
                      "children": [ ]
                    }
                  ]
                }
              ]
            },
            {
              "id": '21',
              "name": "美国",
              "capacity": 1000,
              "used": 0,
              "free": 1000,
              "children": [
                {
                  "id": '52',
                  "name": "杭州",
                  "capacity": 100,
                  "used": 0,
                  "free": 100,
                  "children": [ ]
                }
              ]
            }
          ]
        }
      ]
    }
  ]
};
export default class Container extends React.Component {
  state = {
    searchKeys: [
      { label: '单选', key: 'radio' },
      { label: '多选', key: 'checkbox' },
      { label: '输入框', key: 'input' },
      { label: '时间选择器', key: 'rangePicker' },
      { label: '数字输入框', key: 'number' },
      { label: '树', key: 'tree' },
      { label: '级联选择同步', key: 'cascader_ysnc' },
      { label: '级联选择异步', key: 'cascader_aysnc' },
      { label: '区块', key: 'section' }
    ],
    searchValues: {
      'radio': { type: 'radio', list: [{ label: 'A1', value: 'a1' }, { label: 'A2', value: 'a2' }] },
      'checkbox': { type: 'checkbox', list: [{ label: 'B1', value: 'b1' }, { label: 'B2', value: 'b2' }] },
      'input': { type: 'input', placeholder: '请输入，按回车键确定' },
      'rangePicker': { type: 'rangePicker', placeholder: [ '开始时间', '结束时间' ] },
      'number': { type: 'number', placeholder: '请输入，按回车键确定' },
      'tree': { type: 'tree', data: mockData },
      'cascader_ysnc': { type: 'cascader_ysnc', list: [{ label: 'A1', value: 'a1', children: [{ name: 'children1', id: 1 }, { name: 'children2', id: 2 }] }, { label: 'A2', value: 'a2', children: [{ name: 'children3', id: 3 }, { name: 'children4', id: 4 }] }] },
      'cascader_aysnc': {
        grade: 2,
        type: 'cascader_aysnc',
        list: [{ value: 1, label: "test" }],
        getSecondList: (id) => this.getSecondList(id),
        secondList: [],
        getThirdList: (id) => this.getThirdList(id),
        thirdList: [],
        getFourthList: (id) => this.getFourthList(id),
        fourthList: []
      },
      'section': { type: 'section' }
    }
  };


  async getSecondList(id) {
    const res = await getWithArgs('/api/cloudboot/v1/server-rooms', { page: 1, page_size: 100, idc_id: id });
    if (res.status !== 'success') {
      return;
    }
    const { searchValues } = this.state;
    searchValues.cascader_aysnc.secondList = res.content.records.map(it => {
      return {
        label: it.name,
        value: it.id
      };
    });
    this.setState({ searchValues });
  }

  async getThirdList(id) {
    const res = await getWithArgs('/api/cloudboot/v1/server-cabinets', { page: 1, page_size: 100, server_room_id: id });
    if (res.status !== 'success') {
      return;
    }
    const { searchValues } = this.state;
    searchValues.cascader_aysnc.thirdList = res.content.records.map(it => {
      return {
        label: it.number,
        value: it.id
      };
    });
    this.setState({ searchValues });
  }

  async getFourthList(id) {
    const res = await getWithArgs('/api/cloudboot/v1/server-usites', { page: 1, page_size: 100, server_cabinet_id: id });
    if (res.status !== 'success') {
      return;
    }
    const { searchValues } = this.state;
    searchValues.cascader_aysnc.fourthList = res.content.records.map(it => {
      return {
        label: it.number,
        value: it.id
      };
    });
    this.setState({ searchValues });
  }

  onSearch = (data) => {
    console.log(data);
  };

  render() {
    return (
      <SearchForm onSearch={this.onSearch} searchKeys={this.state.searchKeys} searchValues={this.state.searchValues}/>
    );
  }
}
