#!/bin/bash

# PDF矢量图提取工具构建脚本
echo "开始构建PDF矢量图提取工具..."

# 清理之前的构建
echo "清理之前的构建文件..."
rm -rf build/
mkdir -p build

# 下载依赖
echo "下载Go模块依赖..."
go mod tidy
go mod download

# Windows 64位构建
echo "构建Windows 64位版本..."
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1

# 设置静态链接标志
export CGO_LDFLAGS="-static -static-libgcc -static-libstdc++"

# 构建可执行文件
go build -ldflags "-w -s -H=windowsgui -extldflags '-static'" -o build/pdf-extractor-windows-amd64.exe .

if [ $? -eq 0 ]; then
    echo "✅ Windows 64位版本构建成功: build/pdf-extractor-windows-amd64.exe"
    ls -lh build/pdf-extractor-windows-amd64.exe
else
    echo "❌ Windows 64位版本构建失败"
    exit 1
fi

# Linux 64位构建（用于测试）
echo "构建Linux 64位版本（测试用）..."
export GOOS=linux
export GOARCH=amd64

go build -ldflags "-w -s" -o build/pdf-extractor-linux-amd64 .

if [ $? -eq 0 ]; then
    echo "✅ Linux 64位版本构建成功: build/pdf-extractor-linux-amd64"
    ls -lh build/pdf-extractor-linux-amd64
else
    echo "❌ Linux 64位版本构建失败"
fi

# macOS 64位构建
echo "构建macOS 64位版本..."
export GOOS=darwin
export GOARCH=amd64

go build -ldflags "-w -s" -o build/pdf-extractor-macos-amd64 .

if [ $? -eq 0 ]; then
    echo "✅ macOS 64位版本构建成功: build/pdf-extractor-macos-amd64"
    ls -lh build/pdf-extractor-macos-amd64
else
    echo "❌ macOS 64位版本构建失败"
fi

echo ""
echo "构建完成！生成的文件："
ls -la build/

echo ""
echo "主要目标文件（Windows）："
echo "build/pdf-extractor-windows-amd64.exe - 可直接在Windows上运行"