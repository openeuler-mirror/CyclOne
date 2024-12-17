const fs = require('fs');
const pageContainers = fs.readdirSync('app/containers');

/**
 * 获取Container列表
 *
 * @name ContainerList
 * @function
 * @returns {Array}
 */
function containerList() {
  return pageContainers;
}

module.exports = containerList;
