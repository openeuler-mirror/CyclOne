export function resolve(path, obj, key, value) {
  let recursivePath = '';
  return path.split('.').reduce((prev, curr) => {
    if (value) {
      recursivePath += curr + '.';
      if (recursivePath.substr(0, recursivePath.length - 1) === path) {
        prev[curr][key] = value;
      }
    }
    return prev ? prev[curr] : undefined;
  }, obj || {});
}
