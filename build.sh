#!/bin/bash

rm -rf bin/*
versionDir="github.com/shmy/dd-server/pkg/version"
# 获取gitTag
gitTag=$(if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
# 获取gitCommit
gitCommit=$(git log --pretty=format:'%H' -n 1)
# 获取 gitTreeState
gitTreeState=$(if git status|grep -q 'clean';then echo clean; else echo dirty; fi)
# 获取打包时间
buildDate="$(TZ=Asia/Shanghai date +%FT%T%z)"

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState} -X ${versionDir}.buildDate=${buildDate} -s"

# go-bindata -o=util/asset.go -pkg=util public/web_client/...

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/darwin/dd-server -v -ldflags "$ldflags" .

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/dd-server -v -ldflags "$ldflags" .

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/dd-server.exe -v -ldflags "$ldflags" .

# 压缩
upx bin/darwin/dd-server
upx bin/linux/dd-server
upx bin/windows/dd-server.exe

# 配置文件
cp config.prod.yml bin/darwin/config.yml
cp config.prod.yml bin/linux/config.yml
cp config.prod.yml bin/windows/config.yml

# 静态资源
cp -R public bin/darwin/public
cp -R public bin/linux/public
cp -R public bin/windows/public