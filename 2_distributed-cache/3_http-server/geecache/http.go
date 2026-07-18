// 定义geecache包，实现分布式缓存相关功能
package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 默认的基础路径，用于HTTP请求的路由前缀
const defaultBasePath = "/_geecache/"

// HTTPPool 实现了PeerPicker接口，用于管理HTTP节点池（分布式缓存中的节点）
type HTTPPool struct {
	self     string // 当前节点的基础URL，例如 "https://example.net:8000"
	basePath string // 缓存服务的HTTP请求基础路径，默认是defaultBasePath
}

// NewHTTPPool 初始化一个HTTP节点池实例
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,            // 初始化当前节点的URL
		basePath: defaultBasePath, // 使用默认的基础路径
	}
}

// Log 带服务器名称的日志输出，方便区分不同节点的日志
func (p *HTTPPool) Log(format string, v ...interface{}) {
	// 格式化日志，前缀包含当前节点的self标识
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 处理所有的HTTP请求，实现了http.Handler接口
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 检查请求路径是否以basePath为前缀，确保只处理缓存相关的请求
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	// 记录请求的方法和路径到日志
	p.Log("%s %s", r.Method, r.URL.Path)

	// 请求路径格式要求：/<basepath>/<groupname>/<key>
	// 去除basePath前缀后，按"/"分割为两部分（groupname和key）
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	// 如果分割后不是两部分（即缺少groupname或key），返回400错误
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0] // 从分割结果中提取缓存组名称
	key := parts[1]       // 从分割结果中提取缓存键

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
