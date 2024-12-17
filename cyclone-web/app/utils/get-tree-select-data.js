
//机架和机位树结构数据转换
export default function renderTreeData(data, code) {
  if (!data) {
    return null;
  }
  if (!(data instanceof Array)) {
    data = [data];
  }
  // title: 'Node1',
  // value: '0-0',
  // key: '0-0',
  const loop = (el) => {
    el.title = `${el.cabinet_number} 可用机位数: ${el.available_usites_count}`;
    el.value = el.cabinet_number;
    el.key = el.cabinet_number;
    if (el.leaves) {
      el.leaves.forEach(t => {
        t.title = t.usite_number;
        t.value = t.usite_id;
        t.key = t.usite_id;
      });
      el.children = el.leaves;
    }
  };
  data.forEach(loop);
  console.log(data);
  return data;
}
