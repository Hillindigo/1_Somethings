package main

/*
// 示例curl命令，用于测试API服务的响应结果
// 1. 请求存在的key（Tom），预期返回分数630
$ curl "http://localhost:9999/api?key=Tom"
630

// 2. 请求不存在的key（kkk），预期返回错误信息
$ curl "http://localhost:9999/api?key=kkk"
kkk not exist
*/

import (
	"Dcache/5_multi-nodes/geecache" // 导入自定义的分布式缓存包，用于实现节点间通信和数据缓存
	"flag"                          // 用于解析命令行参数
	"fmt"
	"log"      // 用于日志输出
	"net/http" // 用于启动HTTP服务和处理请求
)

// db 模拟本地数据库，存储用户（key）与分数（value）的映射关系
var db = map[string]string{
	"Tom":  "630", // Tom的分数为630
	"Jack": "589", // Jack的分数为589
	"Sam":  "567", // Sam的分数为567
}

// createGroup 创建并返回一个geecache.Group实例（缓存组）
func createGroup() *geecache.Group {
	// 调用geecache.NewGroup创建缓存组，参数分别为：
	// 1. 组名"scores"；2. 缓存最大容量2<<10字节（2KB）；3. 数据加载函数（缓存未命中时从db加载数据）
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key) // 打印日志，标记从"慢数据库"（模拟db）查询
			if v, ok := db[key]; ok {               // 检查db中是否存在该key
				return []byte(v), nil // 存在则返回字节切片形式的值
			}
			// 不存在则返回错误信息
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// startCacheServer 启动缓存节点服务，实现分布式缓存的节点功能
// 参数说明：addr-当前节点地址；addrs-所有节点地址列表；gee-关联的缓存组
func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	// 1. 创建HTTP节点池实例，传入当前节点地址
	peers := geecache.NewHTTPPool(addr)
	// 2. 设置节点池中的所有节点地址（用于一致性哈希选择节点）
	peers.Set(addrs...)
	// 3. 将节点池注册到缓存组，使缓存组能通过节点池选择节点获取数据
	gee.RegisterPeers(peers)
	// 4. 打印当前节点启动日志
	log.Println("geecache is running at", addr)
	// 5. 启动HTTP服务：监听addr中端口部分（如"http://localhost:8001"取"localhost:8001"），用peers处理请求
	// log.Fatal会在服务启动失败时打印错误并退出程序
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

// startAPIServer 启动API服务，提供对外的缓存查询接口（供客户端通过URL参数查询缓存）
// 参数说明：apiAddr-API服务地址；gee-关联的缓存组
func startAPIServer(apiAddr string, gee *geecache.Group) {
	// 1. 注册"/api"路径的请求处理器，使用匿名函数实现http.HandlerFunc
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// 从URL查询参数中获取"key"的值（如"?key=Tom"中的"Tom"）
			key := r.URL.Query().Get("key")
			// 从缓存组中获取key对应的值
			view, err := gee.Get(key)
			if err != nil {
				// 获取失败：返回500错误和错误信息
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// 获取成功：设置响应头为二进制流类型，写入缓存值的字节切片
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	// 2. 打印API服务启动日志
	log.Println("fontend server is running at", apiAddr)
	// 3. 启动API服务：监听apiAddr中端口部分，使用默认的ServeMux处理请求
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

// main 程序入口，负责解析参数、初始化组件并启动服务
func main() {
	// 1. 定义命令行参数：
	// -port：缓存节点服务端口，默认8001；-api：是否启动API服务，默认false
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port") // 绑定port变量到-port参数
	flag.BoolVar(&api, "api", false, "Start a api server?")  // 绑定api变量到-api参数
	flag.Parse()                                             // 解析命令行参数

	// 2. 定义API服务地址和节点地址映射（端口到完整URL）
	apiAddr := "http://localhost:9999" // API服务固定地址
	addrMap := map[int]string{         // 不同端口对应的缓存节点地址
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	// 3. 从addrMap提取所有节点地址，存入addrs切片
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	// 4. 创建缓存组实例
	gee := createGroup()

	// 5. 如果命令行参数-api为true，启动API服务（使用协程，避免阻塞缓存节点服务）
	if api {
		go startAPIServer(apiAddr, gee)
	}

	// 6. 启动缓存节点服务：根据命令行参数-port选择节点地址，传入所有节点地址和缓存组
	startCacheServer(addrMap[port], addrs, gee)
}
