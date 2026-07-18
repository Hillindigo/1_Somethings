package cache

import (
	"Dcache/2_single-node/geecache/cache/lru" // 导入自定义的lru缓存包，用于实现LRU缓存策略
	"sync"                                    // 导入同步包，用于实现并发安全
)

// cache 结构体封装了LRU缓存，并通过互斥锁保证并发安全
type cache struct {
	mu         sync.Mutex // 互斥锁，用于保护缓存操作的并发安全
	lru        *lru.Cache // 指向lru.Cache实例，实际存储缓存数据
	cacheBytes int64      // 缓存允许使用的最大字节数
}

// add 方法用于向缓存中添加键值对，添加时会加锁保证并发安全
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()         // 加锁，防止并发冲突
	defer c.mu.Unlock() // 函数结束时解锁

	// 如果lru实例未初始化，则创建一个新的LRU缓存实例
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	// 调用lru.Cache的Add方法添加键值对
	c.lru.Add(key, value)
}

// get 方法用于从缓存中获取键对应的值，获取时会加锁保证并发安全
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()         // 加锁，防止并发冲突
	defer c.mu.Unlock() // 函数结束时解锁

	// 如果lru实例未初始化，直接返回（此时缓存为空）
	if c.lru == nil {
		return
	}

	// 调用lru.Cache的Get方法获取值，若存在则类型断言为ByteView后返回
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	// 键不存在时返回空值和false
	return
}
