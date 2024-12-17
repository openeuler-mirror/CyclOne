#!/bin/sh
# 文件打包脚本
# date：2013/2/27


echo '************************************************************************************************'
echo '*  这个是包管理脚本,主要当前项目打包                                                                *'
echo '*  参数1(必选):包名称,一般是仓库名+分支名(master分支后缀) 例如:pkg-manager-develop,pkg-manager       *'
echo '*  参数2(必选):包版本,一般是{number}.{number}.{number}的命名方式,例如 1.0.0，2.0.0                  *'
echo '*  参数3(可选):打包服务其地址,比如 dcospkg.stable.idcos.net, 默认是 dcospkg.stable.idcos.net        *'
echo '************************************************************************************************'

PKG_SERVER="dcospkg.stable.idcos.net"

if [ $# -lt 2 ];then
	echo "??????????????????请输入包名称和包版本号????????????????????"
	exit
fi

git_branch=`echo $1 | sed -E 's/(^[^/]+)?\///g'`

if [ $git_branch == "master" ];then
	BRANCH=""
else
	BRANCH="-"$git_branch
fi

echo $BRANCH

cd pkgcmd

if [ $# -eq 2 ];then

	echo "ant -f build.xml -Dbranch=$BRANCH -Dpackage.version=$2 -Dpackage.server=$PKG_SERVER"
	ant -f build.xml -Dbranch=$BRANCH -Dpackage.version=$2 -Dpackage.server=$PKG_SERVER

elif [ $# -eq 3 ];then 

	echo "ant -f build.xml -Dbranch=$BRANCH -Dpackage.version=$2 -Dpackage.server=$3"
	ant -f build.xml -Dbranch=$BRANCH -Dpackage.version=$2 -Dpackage.server=$3

else
	echo "??????????????????请输入包名称和包版本号????????????????????"
fi


