package GeeCache

import "github.com/linyerun/GeeCache/cache"

type INodeClientGetter interface { // 根据某个分组的key去获取一个发请求的客户端
	GetNodeClient(key string) (nodeClient INodeClient, ok bool)
	SetNodeClient(clientBaseUrls ...string)
}

type INodeClient interface { // 有了节点客户端, 我们就可以发请求去获取数据了
	Get(group string, key string) ([]byte, error)
	Post(group string, key string, value []byte) error
	Delete(group string, key string) error
	GetBaseUrl() string
}

type IServer interface {
	Run() error
}

type IGroup interface {
	Get(key string) (cache.IByteView, error)
	AddOrUpdate(key string, value cache.IByteView) (err error)
	Delete(key string) (err error)
}

type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 效仿HandleFunc的设计
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}
