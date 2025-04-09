@echo off
setlocal enabledelayedexpansion

REM Proto 根目录
set PROTO_ROOT=.
set OUT_DIR=.

REM 编译参数
set OPTIONS=--proto_path=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative

echo 🚀 Generating gRPC code...

REM 编译指定目录下的所有 proto 文件
protoc %OPTIONS% "user\*.proto"

REM 继续添加其它模块
REM protoc %OPTIONS% "product\*.proto"

echo ✅ All proto files compiled!
pause
