# PDF矢量图提取工具 Makefile

.PHONY: all clean windows linux macos test deps

# 默认目标
all: windows

# 清理构建文件
clean:
	rm -rf build/

# 创建构建目录
build-dir:
	mkdir -p build

# 下载依赖
deps:
	go mod tidy
	go mod download

# Windows静态构建
windows: build-dir deps
	@echo "构建Windows 64位静态版本..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
	go build -ldflags "-w -s -H=windowsgui -extldflags '-static'" \
	-o build/pdf-extractor-windows-amd64.exe .
	@echo "✅ Windows版本构建完成: build/pdf-extractor-windows-amd64.exe"

# Linux构建
linux: build-dir deps
	@echo "构建Linux 64位版本..."
	GOOS=linux GOARCH=amd64 \
	go build -ldflags "-w -s" \
	-o build/pdf-extractor-linux-amd64 .
	@echo "✅ Linux版本构建完成: build/pdf-extractor-linux-amd64"

# macOS构建
macos: build-dir deps
	@echo "构建macOS 64位版本..."
	GOOS=darwin GOARCH=amd64 \
	go build -ldflags "-w -s" \
	-o build/pdf-extractor-macos-amd64 .
	@echo "✅ macOS版本构建完成: build/pdf-extractor-macos-amd64"

# 构建所有平台
all-platforms: windows linux macos

# 运行测试
test:
	go test -v ./...

# 运行当前平台版本（用于开发测试）
run:
	go run .

# 检查Go模块
check:
	go mod verify
	go vet ./...

# 显示构建信息
info:
	@echo "Go版本: $(shell go version)"
	@echo "构建目标: windows/amd64 (主要目标)"
	@echo "输出文件: build/pdf-extractor-windows-amd64.exe"
	@echo ""
	@echo "可用命令:"
	@echo "  make windows    - 构建Windows版本"
	@echo "  make linux      - 构建Linux版本"
	@echo "  make macos      - 构建macOS版本"
	@echo "  make all-platforms - 构建所有平台"
	@echo "  make clean      - 清理构建文件"
	@echo "  make test       - 运行测试"
	@echo "  make run        - 运行开发版本"