.PHONY: run build tidy clean

# 启动服务
run:
	go run cmd/server/main.go

# 编译构建
build:
	go build -o bin/blog cmd/server/main.go

# 整理依赖
tidy:
	go mod tidy

# 清理构建产物
clean:
	rm -rf bin/
