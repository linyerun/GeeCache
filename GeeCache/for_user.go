package GeeCache

import (
	"fmt"
	"github.com/linyerun/GeeCache/cache"
	"github.com/linyerun/GeeCache/consistent_hash"
	"github.com/linyerun/GeeCache/single_fighting"
	"github.com/linyerun/GeeCache/utils"
	"strings"
)

type SetNodesURL struct{}

func (_ SetNodesURL) Set(clientBaseUrls ...string) {
	clientsGetter.SetNodeClient(clientBaseUrls...)
}

// NewHttpClientGetter 使用一致性Hash算法开启分布式节点
func NewHttpClientGetter(replicas int, fn consistent_hash.HashFunc) SetNodesURL {
	clientsGetter = &httpClientsGetter{
		nodes:       consistent_hash.NewConsistentHash(replicas, fn),
		httpClients: make(map[string]INodeClient),
	}
	return SetNodesURL{}
}

// NewHttpServer 使用Http服务
func NewHttpServer(ip string, port uint, basePath string) IServer {
	hostBaseUrl = fmt.Sprintf("%s:%d/", ip, port)
	parts := strings.Split(basePath, "/")
	basePath = ""
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		basePath += part + "/"
	}
	globalBasePath = basePath
	return &HttpServer{
		port:     port,
		ip:       ip,
		basePath: basePath,
	}
}

// NewGroupByGetter 使用Getter接口创建Group
func NewGroupByGetter(groupName string, maxCacheBytes int64, getter Getter, onEvicted func(key string, value cache.IByteView)) (g IGroup) {
	if getter == nil {
		utils.Logger().Errorln("getter can not be nil")
		panic("getter is nil!")
	}
	groupsRwMutex.Lock()
	defer groupsRwMutex.Unlock()
	g = &Group{
		name:           groupName,
		getter:         getter,
		safetyCache:    cache.NewSafetyCache(maxCacheBytes, onEvicted),
		singleFighting: single_fighting.NewSingleFighting(),
	}
	groups[groupName] = g
	return
}

// NewGroupByGetterFunc 使用GetterFunc函数创建Group
func NewGroupByGetterFunc(groupName string, cacheBytes int64, f GetterFunc, onEvicted func(key string, value cache.IByteView)) (g IGroup) {
	return NewGroupByGetter(groupName, cacheBytes, f, onEvicted)
}
