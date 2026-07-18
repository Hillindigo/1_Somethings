// 定义geecache包，用于实现分布式缓存相关功能
package geecache

// PeerPicker 是一个接口，必须实现该接口以定位持有特定key的节点（peer）
// 作用：在分布式缓存中，根据key选择对应的节点来获取数据
type PeerPicker interface {
	// PickPeer 根据key选择对应的节点，返回该节点的PeerGetter接口和是否成功的标志
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 是一个接口，必须由节点（peer）实现，用于从节点获取数据
// 作用：定义节点提供的缓存查询能力，供其他节点调用
type PeerGetter interface {
	// Get 从指定的group中获取key对应的数据，返回字节切片和可能的错误
	Get(group string, key string) ([]byte, error)
}
