export default (
  data = [],
  keyMap = {
    valueKey: 'id',
    labelKey: 'name'
  }
) => {
  return data.map(it => {
    return {
      label: it[keyMap.labelKey],
      value: it[keyMap.valueKey]
    };
  });
};
