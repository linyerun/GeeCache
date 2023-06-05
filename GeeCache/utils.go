package GeeCache

// GetGroup 获取分组
func GetGroup(groupName string) (g IGroup, ok bool) {
	groupsRwMutex.RLock()
	defer groupsRwMutex.RUnlock()
	g, ok = groups[groupName]
	return
}

// 改进一下path的前后下划线问题
func decodeBasePath(basePath string) string { // 1:/, 2:/xx/xx/
	if len(basePath) == 0 {
		return "/"
	}
	if basePath[0] != '/' {
		basePath = "/" + basePath
	}
	if basePath[len(basePath)-1] != '/' {
		basePath = basePath + "/"
	}
	return basePath
}
