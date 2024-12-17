#!/bin/bash
#set -x
#******************************************************************************
# @file    : hf-pub.sh
# @author  : wangyubin
# @date    : 2016-04-21 13:42:34
# 
# @brief   : publish script for hf env
# history  : init
#******************************************************************************

# git tag
TAG=`date "+%Y%m%d%H%M%S"`
TAG_CONTENT=""
TAG_FILE=""
PKG_NAME=master
PKG_SERVER=dcospkg.stable.idcos.net


function usage() {
    echo "hf-pub.sh -c TAG_CONTENT OR -F TAG_FILE"
    # echo "        -p PKG_NAME  -s <PKG_SERVER>"
    # echo "          -c TAG_CONTENT OR -F TAG_FILE"
    # echo "    -p PKG_NAME:    包名称       **必选参数**"
    # echo "    -s PKG_SERVER:  打包服务器地址  可选参数， 默认 55.3.15.142"
    echo "    -c TAG_CONTENT: tag 内容 指定 -c 后不能指定 -F"
    echo "    -F TAG_FILE:    tag 内容描述文件 指定 -F 后不能指定 -c"
    echo ""
    echo ""
}

while getopts "c:p:s:F:" arg
do
    case $arg in
        p)
            PKG_NAME=$OPTARG
            ;;
        s)
            PKG_SERVER=$OPTARG
            ;;
        c)
            TAG_CONTENT=$OPTARG
            ;;
        F)
            TAG_FILE=$OPTARG
            ;;
        ?)
            echo "unkonw argument"
            usage
            exit 1
            ;;
    esac
done

if [ "$TAG_CONTENT" == "" -a "$TAG_FILE" == "" ]; then
    echo "必须指定 -c 或者 -F 参数"
    usage
    exit 1
fi

if [ "$TAG_CONTENT" != "" -a "$TAG_FILE" != "" ]; then
    echo "-c 或者 -F 参数 只能指定一个"
    usage
    exit 1
fi

if [ "$TAG_CONTENT" != "" ]; then
    git tag -m $TAG_CONTENT $TAG
else
    git tag -F $TAG_FILE $TAG
fi

if [ $? -ne 0 ]; then
    exit 1
fi

git tag -l -n20 $TAG

# git push --tag

echo "5秒后开始打包 ..."
sleep 5

sh pkg.sh $PKG_NAME $TAG $PKG_SERVER

