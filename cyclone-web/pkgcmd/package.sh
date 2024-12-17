
CURRENT_DIR=`pwd`
echo "执行打包脚本，当前工作目录:$CURRENT_DIR"
cd $CURRENT_DIR/temp
zip -r -q $1 * 