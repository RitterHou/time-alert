#!/usr/bin/env bash

# 生成包含了图标的coff文件
rsrc -ico icon.ico -o TimeAlert.syso
# 执行编译操作，编译时会自动连接之前生成的coff文件
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui -linkmode internal -s -w" -o TimeAlert.exe