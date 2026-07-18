// consistenthash 包实现了一致性哈希算法，用于分布式系统中节点的负载均衡
package consistenthash

import (
	"hash/crc32" // 导入crc32哈希函数库，用于计算哈希值
	"sort"       // 导入排序库，用于对哈希环上的节点进行排序
	"strconv"    // 导入字符串转换库，用于生成虚拟节点的标识
)

// Hash 是一个函数类型，定义了将字节数据映射为uint32哈希值的方法
type Hash func(data []byte) uint32

// Map 结构体表示一致性哈希环，包含所有哈希后的键和节点映射关系
type Map struct {
	hash     Hash           // 哈希函数，用于计算键和节点的哈希值
	replicas int            // 每个真实节点对应的虚拟节点数量
	keys     []int          // 存储所有虚拟节点的哈希值，按升序排列（哈希环）
	hashMap  map[int]string // 虚拟节点哈希值到真实节点的映射
}

// New 创建一个新的一致性哈希Map实例
// 参数replicas是每个真实节点的虚拟节点数量，fn是自定义哈希函数（可为nil，默认使用crc32）
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,             // 初始化虚拟节点数量
		hash:     fn,                   // 初始化哈希函数
		hashMap:  make(map[int]string), // 初始化虚拟节点到真实节点的映射表
	}
	// 如果未指定哈希函数，默认使用crc32.ChecksumIEEE
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 向哈希环中添加一个或多个真实节点（键）
func (m *Map) Add(keys ...string) {
	// 遍历每个真实节点
	for _, key := range keys {
		// 为每个真实节点创建replicas个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 生成虚拟节点的标识（格式：i+key，如"0server1"），计算其哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// 将虚拟节点的哈希值添加到哈希环（keys切片）
			m.keys = append(m.keys, hash)
			// 记录虚拟节点哈希值到真实节点的映射
			m.hashMap[hash] = key
		}
	}
	// 对哈希环上的虚拟节点哈希值进行升序排序，形成环形结构
	sort.Ints(m.keys)
}

// Get 根据输入的键，查找哈希环上最近的虚拟节点对应的真实节点
func (m *Map) Get(key string) string {
	// 如果哈希环为空（没有节点），返回空字符串
	if len(m.keys) == 0 {
		return ""
	}

	// 计算输入键的哈希值
	hash := int(m.hash([]byte(key)))
	// 使用二分查找，找到哈希环上第一个大于等于当前键哈希值的虚拟节点索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 如果索引等于哈希环长度（即超出环的末尾），则取第一个节点（形成环）
	// 否则取索引对应的虚拟节点，通过hashMap找到真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
