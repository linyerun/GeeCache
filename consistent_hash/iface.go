package consistent_hash

// HashFunc 哈希算法
type HashFunc func(data []byte) uint32

type IConsistentHash interface {
	Add(keys ...string)
	Get(key string) (string, bool)
}
