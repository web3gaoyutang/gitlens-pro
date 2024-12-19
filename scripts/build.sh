#!/bin/bash

# 编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o bin/activate.exe

# 编译 Linux 版本
GOOS=linux GOARCH=amd64 go build -o bin/activate

# 编译 MacOS 版本
GOOS=darwin GOARCH=amd64 go build -o bin/activate_mac

echo "构建完成!" 