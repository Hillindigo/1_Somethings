// 定义lru包，用于实现LRU缓存功能
package lru

// 导入标准库的双向链表包，用于维护缓存元素的访问顺序
import "container/list"

// Cache 是LRU缓存的结构体，不支持并发访问
type Cache struct {
	maxBytes  int64                         // 缓存允许使用的最大字节数
	nbytes    int64                         // 当前缓存已使用的字节数
	ll        *list.List                    // 双向链表，用于记录元素的访问顺序（最近使用的在头部）
	cache     map[string]*list.Element      // 哈希表，用于快速查找元素（键为字符串，值为链表节点指针）
	OnEvicted func(key string, value Value) // 可选回调函数，当元素被淘汰时执行
}

// entry 是双向链表中存储的元素结构，包含键和值
type entry struct {
	key   string // 缓存的键
	value Value  // 缓存的值（实现Value接口）
}

// Value 接口定义了缓存值需要实现的方法，用于计算值占用的字节数
type Value interface {
	Len() int // 返回值所占用的字节数
}

// New 是Cache的构造函数，用于创建一个新的LRU缓存实例
// 参数maxBytes为缓存最大容量（字节），onEvicted为元素被淘汰时的回调函数（可选）
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,                       // 初始化最大字节数
		ll:        list.New(),                     // 初始化双向链表
		cache:     make(map[string]*list.Element), // 初始化哈希表
		OnEvicted: onEvicted,                      // 初始化回调函数
	}
}

// Add 方法用于向缓存中添加或更新一个键值对
func (c *Cache) Add(key string, value Value) {
	// 检查键是否已存在于缓存中
	if ele, ok := c.cache[key]; ok {
		// 键已存在：将对应的链表节点移到头部（标记为最近使用）
		c.ll.MoveToFront(ele)
		// 将链表节点的值转换为entry类型
		kv := ele.Value.(*entry)
		// 更新当前使用的字节数（新值长度 - 旧值长度）
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		// 更新节点中的值
		kv.value = value
	} else {
		// 键不存在：在链表头部插入新节点（新元素为最近使用）
		ele := c.ll.PushFront(&entry{key, value})
		// 将新节点存入哈希表，建立键到节点的映射
		c.cache[key] = ele
		// 累加当前使用的字节数（键的长度 + 值的长度）
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 当缓存使用字节数超过最大限制时，循环淘汰最旧的元素
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get 方法用于查询缓存中指定键的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 从哈希表中查找键对应的链表节点
	if ele, ok := c.cache[key]; ok {
		// 键存在：将节点移到链表头部（标记为最近使用）
		c.ll.MoveToFront(ele)
		// 提取节点中的值并返回
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	// 键不存在：返回空值和false
	return
}

// RemoveOldest 方法用于淘汰缓存中最旧的元素（最久未使用的元素）
func (c *Cache) RemoveOldest() {
	// 获取链表尾部的节点（最旧的元素）
	ele := c.ll.Back()
	if ele != nil {
		// 从链表中删除该节点
		c.ll.Remove(ele)
		// 将节点的值转换为entry类型
		kv := ele.Value.(*entry)
		// 从哈希表中删除该键的映射
		delete(c.cache, kv.key)
		// 减少当前使用的字节数（减去键和值的长度）
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 如果设置了回调函数，执行回调（传入被淘汰的键和值）
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len 方法用于返回缓存中的元素数量
func (c *Cache) Len() int {
	// 返回双向链表的长度（即缓存元素的数量）
	return c.ll.Len()
}
