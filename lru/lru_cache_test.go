package lru

import (
	"fmt"
	"testing"
)

type String string

var _ Value = (String)("")

func (s String) Len() int {
	return len(s)
}

// go test 即可把 lru上的测试用例全跑了

func TestLruCache(t *testing.T) {
	// 创建Cache
	cache := NewLruCache(40, func(key string, value Value) {
		fmt.Printf("the value of key = %v in this cache will be eliminated.\n", key)
	})

	_ = cache.Add("随风", String("with the wind"))
	_ = cache.Add("磐石", String("store"))
	_ = cache.Add("林叶润", String("Ernie"))

	if value, ok := cache.Get("随风"); ok {
		fmt.Printf("%v\n", value)
	}

	_ = cache.Add("林叶润", String("Golang"))

	if value, ok := cache.Get("林叶润"); ok {
		fmt.Printf("%v\n", value)
	}

	fmt.Println("size =", cache.Size())
}
