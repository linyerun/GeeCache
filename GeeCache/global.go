package GeeCache

import "sync"

var (
	clientsGetter  INodeClientGetter
	groupsRwMutex  sync.RWMutex
	groups         = make(map[string]IGroup)
	hostBaseUrl    string // 格式: ip:port/
	globalBasePath string // 格式: xxx/yyy/ 、空
)

const (
	binaryContentType = "application/octet-stream"
)
