// 定义geecache包，实现缓存相关功能
package geecache

import (
	"fmt"
	"log"
	"sync"
)

// Group 表示一个缓存命名空间，包含该命名空间下的数据及关联的数据加载逻辑
type Group struct {
	name      string // 缓存组的名称，用于唯一标识一个缓存组
	getter    Getter // 数据加载器，当缓存未命中时负责加载数据
	mainCache cache  // 主缓存实例，用于实际存储缓存数据
}

// Getter 是一个接口，定义了根据key加载数据的方法
type Getter interface {
	Get(key string) ([]byte, error) // 输入key，返回对应的字节数据和可能的错误
}

// GetterFunc 是一个函数类型，实现了Getter接口（适配器模式，让普通函数可以作为Getter使用）
type GetterFunc func(key string) ([]byte, error)

// Get 实现Getter接口的Get方法，调用函数本身
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key) // 直接调用函数f处理key
}

var (
	mu     sync.RWMutex              // 读写锁，保护groups映射的并发安全（读多写少场景优化）
	groups = make(map[string]*Group) // 存储所有缓存组的映射，key为组名，value为Group实例
)

// NewGroup 创建一个新的缓存组实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil { // 数据加载器不能为空，否则触发panic
		panic("nil Getter")
	}
	mu.Lock()         // 加写锁，保证创建组时的并发安全
	defer mu.Unlock() // 函数结束时自动解锁
	// 初始化Group实例
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes}, // 初始化主缓存，设置最大字节数
	}
	groups[name] = g // 将新创建的组添加到全局映射中
	return g
}

// GetGroup 根据名称获取已创建的缓存组，若不存在则返回nil
func GetGroup(name string) *Group {
	mu.RLock()        // 加读锁，允许多个读操作并发
	g := groups[name] // 从映射中获取组
	mu.RUnlock()      // 读操作完成后解锁
	return g
}

// Get 从缓存中获取key对应的value
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" { // 检查key是否为空，为空则返回错误
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 先从主缓存中查找key
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit") // 缓存命中，打印日志
		return v, nil
	}

	// 缓存未命中，调用load方法加载数据
	return g.load(key)
}

// load 负责加载key对应的数据（目前仅实现本地加载，可扩展为分布式加载）
func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key) // 调用本地加载方法
}

// getLocally 本地加载数据：通过getter获取数据，并将数据存入缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用getter加载原始数据（字节切片）
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err // 加载失败，返回错误
	}
	// 将原始数据转换为ByteView（保证数据不可变）
	value := ByteView{b: cloneBytes(bytes)}
	// 将加载的数据存入缓存，方便后续查询
	g.populateCache(key, value)
	return value, nil
}

// populateCache 将key和对应的value存入主缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value) // 调用主缓存的add方法添加数据
}
