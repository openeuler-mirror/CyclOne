import React from 'react';
import FormGenerator from './FormGenerator';
import { post, put, del, get } from './common/xFetch2';
export default class CustomForm extends React.Component {

  componentDidMount() {
    if (this.props.getDetail) {
      const self = this;
      async function getDetail() {
        try {
          const ret = await get(`${self.props.detailAPI}`);
          if (ret.status === 'success') {
            const data = ret.item;
            if (self.props.jsonParse && data[self.props.jsonParse]) {
              const jsonParseData = JSON.parse(data[self.props.jsonParse]);
              self.props.parseData.map(key => {
                data[key] = jsonParseData[key];
              });
            }
            self.setState({
              initialValue: data
            });
          }
        } catch (err) {
          console.log(err);
        }
      }
      getDetail();
    }
  }
  constructor(props) {
    super(props);
    this.state = {
      initialValue: {}
    };
  }


  render() {
    let $content = this.renderForm();
    return (
      <div className='custom-form'>
        {$content}
      </div>
    );
  }

  renderForm = () => {
    const props = this.props;
    const initialValue = props.initialValue || this.state.initialValue;
    return (
      <div className='form'>
        <FormGenerator
          schema={props.schema}
          initialValue={initialValue}
          showSubmit={false}
        />
      </div>
    );
  }

}
CustomForm.defaultProps = {
  type: 'form'
};
