rm -rf bin/*

# go-bindata -o=util/asset.go -pkg=util public/web_client/...

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/darwin/dd-server -v -ldflags '-w -s' main.go

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux/dd-server -v -ldflags '-w -s' main.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows/dd-server.exe -v -ldflags '-w -s' main.go

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