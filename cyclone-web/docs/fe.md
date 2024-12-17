## 目前项目中前端开发常有的一些问题


### 关于数据流转

目前看到代码仓库里边的数据流转应该是有三种模式：

1. 通过 redux-thunk dispatch 函数类型的 action 
2. 通过 redux-saga  通过 generate 来 compose action
3. 通过 componentDidMount hook 加载数据放在 state 中

第一点想说的大家如果对 1, 2 两种模式比较熟悉，尽量有这两种方案，优先使用 saga 方案，saga 带来能够在复杂场景中展现极大的灵活性

但是如果不是很理解1，2 两种方案，那请用第三种，不理解的情况下写出来的 redux + saga 会显得极其凌乱，很难维护，目前其实我也不了解大家的理解情况，有需要可以专门一起来分享一下，集中提问，集中回答


第二，如果用 state 来存储页面状态和数据的，一定注意以下这点：

// 错误的写法
```js
 this.setState({
    pageNo: pageNo
 })

 getDataService(this.state.pageNo).then(ret=>{...})
```

// 正确的写法

```js
 this.setState({pageNo}, () => {
   // set state 成功过后的回调内请求数据
   getDataService(this.state.pageNo).then(ret=>{...})
 })
```


### 关于 table 

table 渲染的大家都能理解，问题出在大家对 分页和搜索的数据交互上：

1. 分页数据和搜索 query 数据一定要存在 store 或者 state 中
2. 在获取数据的时候直接从 store 或者 state 中拿到分页数据和搜索 query 去获取数据
3. 搜索或者排序的时候要重置 page 为 1. 不然会出现明显的数据交互 bug

另外在 table 的 action 上，目前的方案有：

1. 可通过 components/TableControlCell 来进行渲染，当多余 三个以上的操作是可以 dropdown more
2. 每个 action 的触发的方案有几种：
    - 通过 Popup 弹出一个脱离当前 ReactRoot 的模态框
    - 每个 action 直接 render Modal
    - 建议用 Popup，这种模式可复用性更高 
3. table action 触发修改完成过后要记得 reload table，不然数据不对


### 关于 表单

- 简单的表单都可以通过 components/FormGenerator 来实现配置化动态生成
- 用 FormGenerator 的优势是可以做依赖管理，具体可以参考应用详情的各个表单
- 如果有弹出模态框，要求用户输入一些表单项，然后保存的这种类型，可以不用重复开发，可以利用 FormGenerator/advance-form 来实现

