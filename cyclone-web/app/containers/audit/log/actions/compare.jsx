import React from 'react';
import Popup from 'components/popup/draw';

export default function action(options) {

  const _old = options.records.source;
  const _new = options.records.destination;
  let column = [];

  if (_old !== 'null') {
    const oldJson = JSON.parse(_old);
    Object.keys(oldJson).forEach((key, index) => {
      if (!column[index]) {
        column[index] = {};
      }
      column[index].key = key;
      column[index].old_value = oldJson[key];
    });
  }
  if (_new !== 'null') {
    const newJson = JSON.parse(_new);
    Object.keys(newJson).forEach((key, index) => {
      if (!column[index]) {
        column[index] = {};
      }
      column[index].key = key;
      column[index].new_value = newJson[key];
    });
  }

  Popup.open({
    title: '对比详情',
    onCancel: () => {
      Popup.close();
    },
    content: (
      <div className='log_compare_table'>
        <table>
          <thead>
            <tr>
              <th>key</th>
              <th>原值</th>
              <th>现值</th>
            </tr>
          </thead>
          <tbody>
            {
            column.map(it => <tr className={it.old_value !== it.new_value ? 'diff_tr' : ''}>
              <td>{it.key}</td>
              <td>{it.old_value}</td>
              <td>{it.new_value}</td>
            </tr>)
          }
          </tbody>
        </table>
      </div>
    )
  });
}
