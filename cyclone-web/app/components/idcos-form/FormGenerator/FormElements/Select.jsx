import React from 'react';
import {
  Cascader,
  TreeSelect,
  Form,
  Input,
  Select,
  Button,
  Checkbox,
  Radio,
  Tooltip,
  Icon
} from 'antd';
import fetch from 'common/xFetch';
const FormItem = Form.Item;
const TreeNode = TreeSelect.TreeNode;
const Option = Select.Option;
import T from 'prop-types';

class MSelect extends React.Component {
  state = {
    options: this.props.options || []
  };
  static propTypes= {
    label: T.string,
    labelColSpan: T.number,
    wrapperColSpan: T.number,
    help: T.string,
    tip: T.string,
    value: T.string,
    rules: T.array,
    hidden: T.bool,
    form: T.object.isRequired,
    path: T.string.isRequired,
    dataSource: T.string,
    cascader: T.bool,
    tree: T.bool,
    multiple: T.bool,
    mode: T.string
  };
  static defaultProps = {
    label: '下拉选择',
    labelColSpan: 4,
    wrapperColSpan: 18,
    placeholder: '请选择',
    hidden: false,
    cascader: false,
    multiple: false
  };
  componentWillReceiveProps(nextProps) {
    const dataSource = nextProps.dataSource;

    if (dataSource) {
      if (dataSource !== this.props.dataSource) {
        this.loadAsyncOptions(dataSource);
      }
    } else {
      this.setState({
        options: nextProps.options || this.state.options
      });
    }
  }
  componentDidMount() {
    if (this.props.dataSource) {
      this.loadAsyncOptions(this.props.dataSource);
    }
  }
  loadAsyncOptions(dataSource) {
    const f = a => a;
    const transform = this.props.transform || f;
    fetch(dataSource).then(
      result => {
        const data = transform(result);
        this.setState({
          options: data.options || this.state.options
        });
      },
      err => {
        console.log('err', err);
      }
    );
  }
  render() {
    const props = this.props;

    const itemProps = {
      label: props.label,
      className: props.hidden ? 'hide' : '',
      labelCol: { span: props.labelColSpan },
      wrapperCol: { span: props.wrapperColSpan }
    };

    if (props.help) {
      itemProps.help = props.help;
    }

    if (props.tip) {
      itemProps.label = (
        <span>
          {props.label}
          <Tooltip title={props.tip}>
            <Icon type='question-circle-o' />
          </Tooltip>
        </span>
      );
    }

    return <FormItem {...itemProps}>{this.renderSelect()}</FormItem>;
  }

  renderSelect() {
    const props = this.props;
    const { getFieldProps } = this.props.form;

    const fieldProps = getFieldProps(props.path, {
      rules: props.rules,
      initialValue: props.initialValue,
      value: null,
      onChange() {
        props.onFieldChange({
          path: props.path,
          args: arguments
        });
      }
    });

    if (props.cascader) {
      return (
        <Cascader
          {...fieldProps}
          placeholder={props.placeholder}
          style={{ width: '100%' }}
          options={this.state.options}
          disabled={props.disabled}
        />
      );
    }

    if (props.tree) {
      return this.renderTree(fieldProps);
    }
    if(props.tag){
        return (
            <Select
                {...fieldProps}
                multiple={props.multiple}
                placeholder={props.placeholder}
                disabled={props.disabled}
                showSearch={true}
                mode={props.mode || null}
                optionFilterProp='children'
                style={{ width: '100%' }}
            >
                {this.state.options.map(r => {
                    return (
                        <Option value={r} key={r}>
                            {r}
                        </Option>
                    );
                })}
            </Select>
        );
    }

    return (
      <Select
        {...fieldProps}
        multiple={props.multiple}
        placeholder={props.placeholder}
        disabled={props.disabled}
        showSearch={true}
        mode={props.mode || null}
        optionFilterProp='children'
        filterOption={(input, option) => option.props.children.toLowerCase().indexOf(input.toLowerCase()) >= 0}
        style={{ width: '100%' }}
      >
        {this.state.options.map(r => {
          return (
            <Option value={r.value} key={r.value + r.label + props.id}>
              {r.label}
            </Option>
          );
        })}
      </Select>
    );
  }

  renderTree(fieldProps) {
    const props = this.props;
    const tree = this.state.options;
    return (
      <TreeSelect
        {...fieldProps}
        placeholder={props.placeholder}
        style={{ width: '100%' }}
        disabled={props.disabled}
        allowClear={props.allowClear}
        treeData={this.renderTreeData(tree)}
        treeDefaultExpandAll
      >
        {this.renderTreeNodes(tree)}
      </TreeSelect>
    );
  }

  renderTreeData(data) {
    if (!(data instanceof Array)) {
      data = [data];
    }

    // label: 'Node1',
    // value: '0-0',
    // key: '0-0',
    const loop = el => {
      el.label = el.name;
      el.value = el.id;
      el.key = el.id;
      if (el.children) {
        el.children.forEach(loop);
      }
    };
    data.forEach(loop);
    return data;
  }

  renderTreeNodes(data) {
    if (!(data instanceof Array)) {
      data = [data];
    }
    return data.map(item => {
      const $title = <span> {item.name} </span>;

      if (item.children && item.children.length > 0) {
        return (
          <TreeNode key={item.id} value={item.id} title={$title}>
            {this.renderTreeNodes(item.children)}
          </TreeNode>
        );
      }
      return <TreeNode key={item.id} value={item.id} title={$title} />;
    });
  }
}

export default MSelect;
