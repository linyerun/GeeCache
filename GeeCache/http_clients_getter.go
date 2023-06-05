package GeeCache

import (
	"github.com/linyerun/GeeCache/consistent_hash"
)

type httpClientsGetter struct {
	nodes       consistent_hash.IConsistentHash
	httpClients map[string]INodeClient
}

func (h *httpClientsGetter) GetNodeClient(key string) (nodeClient INodeClient, ok bool) {
	url, ok := h.nodes.Get(key)
	if !ok {
		return nil, false
	}
	return h.httpClients[url], true
}

func (h *httpClientsGetter) SetNodeClient(clientBaseUrls ...string) {
	h.nodes.Add(clientBaseUrls...)       // 搞出虚拟节点, 放到哈希环中
	for _, url := range clientBaseUrls { // 创建对应可以发请求的客户端
		h.httpClients[url] = NewHttpClient(url)
	}
}
