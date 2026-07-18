package cache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

// 模拟数据库，存储用户分数信息
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 测试Getter接口的实现
func TestGetter(t *testing.T) {
	// 定义一个GetterFunc类型的变量f，实现Getter接口
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	// 期望的返回值
	expect := []byte("key")
	// 调用f的Get方法，验证返回值是否符合预期
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatal("callback failed") // 断言失败则触发致命错误
	}
}

// 测试Group的Get方法，验证缓存命中和未命中的逻辑
func TestGet(t *testing.T) {
	// 用于记录每个key的加载次数
	loadCounts := make(map[string]int, len(db))
	// 创建一个名为"scores"的Group，缓存大小为2<<10字节，设置数据加载器
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key) // 模拟从慢数据库查询
			if v, ok := db[key]; ok {
				// 初始化该key的加载次数
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++ // 加载次数加1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	// 遍历db，验证缓存获取逻辑
	for k, v := range db {
		// 第一次获取，应该从加载器加载并缓存
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		// 第二次获取，应该命中缓存，加载次数不应超过1
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	// 验证不存在的key的处理逻辑
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

// 测试GetGroup方法，验证根据名称获取Group的逻辑
func TestGetGroup(t *testing.T) {
	groupName := "scores"
	// 创建一个Group
	NewGroup(groupName, 2<<10, GetterFunc(
		func(key string) (bytes []byte, err error) { return }))
	// 验证能正确获取到该Group
	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("group %s not exist", groupName)
	}

	// 验证获取不存在的Group时返回nil
	if group := GetGroup(groupName + "111"); group != nil {
		t.Fatalf("expect nil, but %s got", group.name)
	}
}
