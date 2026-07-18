// 定义cache包，用于实现缓存相关功能
package cache

// ByteView 用于持有对字节数据的不可变视图（即数据一旦创建就不能被修改）
type ByteView struct {
	b []byte // 存储实际的字节数据
}

// Len 返回当前视图的数据长度（实现了lru.Value接口的Len方法，用于计算缓存占用的字节数）
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 以字节切片的形式返回数据的副本（避免外部直接修改内部的b字段，保证不可变性）
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b) // 调用cloneBytes函数创建副本并返回
}

// String 以字符串形式返回数据（必要时会创建副本，同样保证内部数据不被外部修改）
func (v ByteView) String() string {
	return string(v.b) // 将字节切片转换为字符串返回（字符串在Go中是不可变的）
}

// cloneBytes 用于创建字节切片的副本（核心工具函数，保证数据的不可变性）
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b)) // 创建一个与原切片长度相同的新切片
	copy(c, b)                // 将原切片的数据复制到新切片中
	return c                  // 返回新切片（副本）
}
