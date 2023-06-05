package consistent_hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

func NewConsistentHash(replicas int, fn HashFunc) IConsistentHash {
	if fn == nil {
		fn = crc32.ChecksumIEEE //使用IEEE多项式返回数据的CRC-32校验和。
	}
	if replicas <= 0 {
		replicas = 10
	}
	return &consistentHash{
		hashFunc: fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

type consistentHash struct {
	hashFunc HashFunc       //哈希算法函数
	replicas int            //虚拟节点 = replicas*真实节点
	keys     []int          //哈希环上的虚拟节点的hash值
	hashMap  map[int]string //虚拟节点与真实节点的映射表
}

func (c *consistentHash) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < c.replicas; i++ { //依次计算虚拟节点的hash，并保存到hashMap中
			hash := int(c.hashFunc([]byte(strconv.Itoa(i) + key))) //TODO 这个可能会重复吧！那这个虚拟节点对应的节点就被占用了
			c.keys = append(c.keys, hash)                          //加入hash环中
			c.hashMap[hash] = key                                  //加入集合中
		}
	}
	sort.Ints(c.keys) //对环的虚拟节点的hash值进行排序
}

func (c *consistentHash) Get(key string) (url string, ok bool) {
	if len(c.keys) == 0 {
		return "", false
	}
	hash := int(c.hashFunc([]byte(key)))
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})
	return c.hashMap[c.keys[idx%len(c.keys)]], true //取余的作用是达到切片尽头之后回到原点
}
