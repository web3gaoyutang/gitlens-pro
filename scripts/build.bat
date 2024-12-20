@echo off

:: 创建 bin 目录（如果不存在）
if not exist bin mkdir bin

:: 编译 Windows 版本
set GOOS=windows
set GOARCH=amd64
go build -o bin\activate.exe

:: 编译 Linux 版本
set GOOS=linux
set GOARCH=amd64
go build -o bin\activate

:: 编译 MacOS 版本 (Intel)
set GOOS=darwin
set GOARCH=amd64
go build -o bin\activate_mac_amd64

:: 编译 MacOS 版本 (M1/M2)
set GOOS=darwin
set GOARCH=arm64
go build -o bin\activate_mac_arm64

echo 构建完成！输出目录：.\bin

pause 