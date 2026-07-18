// 定义geecache包，实现分布式缓存核心功能
package geecache

import (
	"fmt"
	"log"
	"sync"
)

// Group 表示一个缓存命名空间，包含该命名空间下的数据、加载逻辑及分布式节点信息
type Group struct {
	name      string     // 缓存组的唯一名称，用于区分不同命名空间
	getter    Getter     // 数据加载器，缓存未命中时负责加载原始数据
	mainCache cache      // 本地主缓存，用于存储当前节点的缓存数据
	peers     PeerPicker // 节点选择器，用于在分布式环境中选择持有目标数据的节点
}

// Getter 接口定义了数据加载的方法，当缓存未命中时会调用该接口从数据源加载数据
type Getter interface {
	Get(key string) ([]byte, error) // 输入key，返回对应的字节数据和可能的错误
}

// GetterFunc 是一个函数类型，通过实现Getter接口，使普通函数可以作为数据加载器
type GetterFunc func(key string) ([]byte, error)

// Get 实现Getter接口的Get方法，直接调用函数本身处理数据加载
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex              // 读写锁，保护全局groups映射的并发安全（读多写少场景优化）
	groups = make(map[string]*Group) // 存储所有缓存组的全局映射，键为组名，值为Group实例
)

// NewGroup 创建一个新的缓存组实例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil { // 数据加载器不能为空，否则触发panic
		panic("nil Getter")
	}
	mu.Lock() // 加写锁，保证创建组时的并发安全
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes}, // 初始化本地缓存，设置最大字节数
	}
	groups[name] = g // 将新创建的组添加到全局映射
	return g
}

// GetGroup 根据名称获取已创建的缓存组，若不存在则返回nil
func GetGroup(name string) *Group {
	mu.RLock()        // 加读锁，允许多个读操作并发执行
	g := groups[name] // 从全局映射中获取缓存组
	mu.RUnlock()      // 读操作完成后解锁
	return g
}

// Get 从缓存中获取key对应的value，优先查本地缓存，未命中则加载数据
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" { // 检查key是否为空，为空则返回错误
		return ByteView{}, fmt.Errorf("key is required")
	}

	// 先从本地主缓存中查找
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit") // 缓存命中，打印日志
		return v, nil
	}

	// 本地缓存未命中，调用load方法加载数据
	return g.load(key)
}

// RegisterPeers 为缓存组注册节点选择器（PeerPicker），用于分布式环境下选择远程节点
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil { // 防止重复注册节点选择器
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// load 负责加载key对应的数据，优先从远程节点获取，失败则本地加载
func (g *Group) load(key string) (value ByteView, err error) {
	// 如果已注册节点选择器，则尝试从远程节点获取数据
	if g.peers != nil {
		// 调用节点选择器选择合适的远程节点
		if peer, ok := g.peers.PickPeer(key); ok {
			// 从选中的远程节点获取数据
			if value, err = g.getFromPeer(peer, key); err == nil {
				return value, nil // 远程获取成功，返回数据
			}
			// 远程获取失败，打印错误日志
			log.Println("[GeeCache] Failed to get from peer", err)
		}
	}

	// 远程获取失败或无节点选择器，调用本地加载方法
	return g.getLocally(key)
}

// populateCache 将key和对应的value存入本地主缓存
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// getLocally 本地加载数据：通过getter从数据源加载数据，并存入本地缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	// 调用数据加载器获取原始数据（字节切片）
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err // 加载失败，返回错误
	}
	// 将原始数据转换为ByteView（保证数据不可变）
	value := ByteView{b: cloneBytes(bytes)}
	// 将加载的数据存入本地缓存，方便后续查询
	g.populateCache(key, value)
	return value, nil
}

// getFromPeer 从远程节点获取数据：通过PeerGetter接口调用远程节点的Get方法
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	// 调用远程节点的Get方法，传入当前缓存组名称和key
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err // 远程获取失败，返回错误
	}
	// 将远程获取的字节数据转换为ByteView并返回
	return ByteView{b: bytes}, nil
}
