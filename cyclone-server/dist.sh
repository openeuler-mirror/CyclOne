#!/bin/sh
export ROOT_DIR=$PWD
export DST_DIR="$ROOT_DIR/dist"
export NOW=$(date "+%Y%m%d%H%M%S")


mkdir -p $DST_DIR
rm -rf $DST_DIR/*
rm -rf $ROOT_DIR/bin/* $ROOT_DIR/pkg/*

echo "============================交叉编译（Linux-AMD64）============================"
GOOS=linux GOARCH=amd64 gbb --debug
cd $ROOT_DIR/bin/linux_amd64
tar -czv -f $DST_DIR/cloudboot-server-linux-amd64-$NOW.tar.gz ./cloudboot-server
tar -czv -f $DST_DIR/cloudboot-agent-linux-amd64-$NOW.tar.gz ./cloudboot-agent
tar -czv -f $DST_DIR/hw-server-linux-amd64-$NOW.tar.gz ./hw-server
rm -rf $ROOT_DIR/bin/* $ROOT_DIR/pkg/*
cd $ROOT_DIR

echo "============================交叉编译（Linux-ARM64）============================"
GOOS=linux GOARCH=arm64 gbb --debug
cd $ROOT_DIR/bin/linux_arm64
tar -czv -f $DST_DIR/cloudboot-agent-linux-arm64-$NOW.tar.gz ./cloudboot-agent
tar -czv -f $DST_DIR/hw-server-linux-arm64-$NOW.tar.gz ./hw-server
rm -rf $ROOT_DIR/bin/* $ROOT_DIR/pkg/*
cd $ROOT_DIR

# echo "============================交叉编译（windows-amd64）============================"
# GOOS=windows GOARCH=amd64 gbb --debug
# cd $ROOT_DIR/bin/windows_amd64
# tar -czv -f $DST_DIR/peagent-windows-amd64-$NOW.tar.gz ./peagent.exe
# tar -czv -f $DST_DIR/peconfig-windows-amd64-$NOW.tar.gz ./peconfig.exe
# tar -czv -f $DST_DIR/winconfig-windows-amd64-$NOW.tar.gz ./winconfig.exe
# tar -czv -f $DST_DIR/win-image-server-windows-amd64-$NOW.tar.gz ./win-image-server.exe
# rm -rf $ROOT_DIR/bin/* $ROOT_DIR/pkg/*
# cd $ROOT_DIR