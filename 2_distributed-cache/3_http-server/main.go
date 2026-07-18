package main

/*
// 示例 curl 命令，用于测试服务
// 访问存在的key
$ curl http://localhost:9999/_geecache/scores/Tom
630

// 访问不存在的key
$ curl http://localhost:9999/_geecache/scores/kkk
kkk not exist
*/

import (
	"Dcache/3_http-server/geecache" // 导入自定义的geecache包
	"fmt"
	"log"
	"net/http"
)

// 模拟数据库，存储用户分数数据
var db = map[string]string{
	"Tom":  "630", // Tom的分数
	"Jack": "589", // Jack的分数
	"Sam":  "567", // Sam的分数
}

func main() {
	// 创建一个名为"scores"的缓存组，最大容量为2<<10字节（2KB）
	// 并设置数据加载函数（当缓存未命中时从db加载数据）
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key) // 打印从"慢数据库"查询的日志
			if v, ok := db[key]; ok {
				return []byte(v), nil // 若key存在，返回对应的值（转换为字节切片）
			}
			// 若key不存在，返回错误信息
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999" // 服务监听的地址
	// 创建HTTP节点池，传入当前节点的地址
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr) // 打印服务启动日志
	// 启动HTTP服务，监听指定地址，使用peers作为处理器
	// log.Fatal会在服务启动失败时打印错误并退出
	log.Fatal(http.ListenAndServe(addr, peers))
}
