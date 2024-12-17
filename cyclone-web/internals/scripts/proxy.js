import fs from 'fs';

process.stdout.write('generate proxy.config');

// 文件是否存在
// 写入文件

const writeFile = (fileName, data) => new Promise((resolve, reject) => {
  fs.writeFile(fileName, data, (error, value) => (error ? reject(error) : resolve(value)));
});

(async function main() {
  var proxyConfigFile = 'proxy.json';
  if (fs.existsSync(proxyConfigFile)) {
    console.log('文件已经存在');
    return;
  }

  var data = {
    "default": {
      "_api": "http://127.0.0.1:8080",
      "api": "http://127.0.0.1:8080",
      "endpoints": [
        "/res/*",
        "/api/*",
        "/user/*",
        "/auth",
        "/auth/*",
        "/login.html",
        "/static/*",
        "/report/*"
      ]
    }
  };

  await writeFile(proxyConfigFile, `${JSON.stringify(data, null, 2)}\n`);

  console.log('创建proxy.json成功');
}());
