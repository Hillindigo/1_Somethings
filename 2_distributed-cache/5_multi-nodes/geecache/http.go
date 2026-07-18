// 定义geecache包，实现分布式缓存的HTTP节点池功能
package geecache

import (
	"Dcache/4_consistent-hash/geecache/consistenthash" // 导入一致性哈希包，用于节点选择
	"fmt"
	"io/ioutil" // 用于读取HTTP响应体
	"log"       // 用于日志输出
	"net/http"  // 用于HTTP服务和请求
	"net/url"   // 用于URL编码
	"strings"   // 用于字符串处理
	"sync"      // 用于并发安全控制
)

// 常量定义
const (
	defaultBasePath = "/_geecache/" // 默认的HTTP请求基础路径
	defaultReplicas = 50            // 一致性哈希中每个真实节点的默认虚拟节点数量
)

// HTTPPool 实现了PeerPicker接口，用于管理HTTP节点池（分布式缓存中的节点集合）
type HTTPPool struct {
	self        string                 // 当前节点的基础URL，例如 "https://example.net:8000"
	basePath    string                 // 缓存服务的HTTP请求基础路径，默认是defaultBasePath
	mu          sync.Mutex             // 互斥锁，保护peers和httpGetters的并发访问安全
	peers       *consistenthash.Map    // 一致性哈希映射，用于根据key选择节点
	httpGetters map[string]*httpGetter // 存储节点地址到httpGetter的映射，key为节点URL
}

// NewHTTPPool 初始化一个HTTP节点池实例
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,            // 初始化当前节点的URL
		basePath: defaultBasePath, // 使用默认的基础路径
	}
}

// Log 输出带服务器名称的日志，方便区分不同节点的日志
func (p *HTTPPool) Log(format string, v ...interface{}) {
	// 格式化日志，前缀包含当前节点的self标识
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 处理所有HTTP请求，实现了http.Handler接口
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 检查请求路径是否以basePath为前缀，确保只处理缓存相关请求
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	// 记录请求的方法和路径到日志
	p.Log("%s %s", r.Method, r.URL.Path)

	// 请求路径格式要求：/<basepath>/<groupname>/<key>
	// 去除basePath前缀后，按"/"分割为两部分（groupname和key）
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	// 如果分割后不是两部分（缺少groupname或key），返回400错误
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0] // 提取缓存组名称
	key := parts[1]       // 提取缓存键

	// 根据组名获取对应的Group实例
	group := GetGroup(groupName)
	if group == nil {
		// 若组不存在，返回404错误
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// 从缓存组中获取key对应的值
	view, err := group.Get(key)
	if err != nil {
		// 若获取失败，返回500错误
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应的Content-Type为二进制流（通用字节数据类型）
	w.Header().Set("Content-Type", "application/octet-stream")
	// 将缓存值的字节切片写入响应体
	w.Write(view.ByteSlice())
}

// Set 更新节点池中的节点列表
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()         // 加锁，保证并发安全
	defer p.mu.Unlock() // 函数结束后解锁

	// 初始化一致性哈希实例，使用默认的虚拟节点数量和哈希函数
	p.peers = consistenthash.New(defaultReplicas, nil)
	// 将传入的节点列表添加到一致性哈希环中
	p.peers.Add(peers...)
	// 初始化httpGetters映射，容量为节点数量
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	// 为每个节点创建对应的httpGetter（用于发送HTTP请求）
	for _, peer := range peers {
		// httpGetter的baseURL为节点地址+基础路径（如"http://localhost:9999/_geecache/"）
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PickPeer 根据key选择对应的节点（实现PeerPicker接口）
func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()         // 加锁，保证并发安全
	defer p.mu.Unlock() // 函数结束后解锁

	// 通过一致性哈希获取key对应的节点，且该节点不是当前节点（p.self）
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer) // 记录选中的节点日志
		// 返回该节点对应的httpGetter（实现了PeerGetter接口）
		return p.httpGetters[peer], true
	}
	// 未找到合适的节点（或选中当前节点），返回nil和false
	return nil, false
}

// 编译期断言：验证HTTPPool是否实现了PeerPicker接口（若未实现则编译报错）
var _ PeerPicker = (*HTTPPool)(nil)

// httpGetter 用于向其他节点发送HTTP请求获取缓存数据，实现PeerGetter接口
type httpGetter struct {
	baseURL string // 基础URL，格式为"节点地址+基础路径"（如"http://localhost:9999/_geecache/"）
}

// Get 向目标节点发送HTTP请求，获取指定group和key的缓存数据（实现PeerGetter接口）
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// 构建请求URL：baseURL + 编码后的group + "/" + 编码后的key
	// url.QueryEscape用于对特殊字符进行编码（如空格、斜杠等）
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)

	// 发送GET请求到目标节点
	res, err := http.Get(u)
	if err != nil {
		return nil, err // 请求失败，返回错误
	}
	defer res.Body.Close() // 确保响应体在函数结束后关闭

	// 检查响应状态码是否为200（成功）
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	// 读取响应体内容（字节切片）
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	// 返回获取到的字节数据
	return bytes, nil
}

// 编译期断言：验证httpGetter是否实现了PeerGetter接口（若未实现则编译报错）
var _ PeerGetter = (*httpGetter)(nil)
