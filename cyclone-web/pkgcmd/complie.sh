CURRENT_DIR=`pwd`
echo "npm run build 开始编译"
npm run build

echo "执行编译脚本，当前工作目录:$CURRENT_DIR"
cp -r build/* temp/www


