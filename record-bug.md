# 请求打到http_server判断前缀有误
- 如果不加 `\`, 而`req.URL.Path`有斜杠开头, 会导致判断前缀不符合
```go
if !strings.HasPrefix(req.URL.Path, "/"+hs.basePath) {
    utils.Logger().Errorln("HttpServer serving unexpected path: " + req.URL.Path)
    http.Error(resp, "HttpServer serving unexpected path: "+req.URL.Path, http.StatusBadRequest)
    return
}
```
# 关于single_fighting的key问题
- 加 `-GET` 和 `-DELETE` 的目的是避免同一个key过来执行GET和DELETE操作共享结果
```go
func (g *Group) Get(key string) (cache.IByteView, error) {
	// 判断缓存中是否有对应的key
	if val, ok := g.safetyCache.Get(key); ok {
		utils.Logger().Println("[GeeCache] Hit key is", key)
		return val, nil
	}

	data, err := g.singleFighting.Do(key+"-GET", func() (any, error) {
		// 判断是否开启了分布式节点模式
		if clientsGetter != nil { // 开启了
			client, _ := clientsGetter.GetNodeClient(key)
			if client.GetBaseUrl() != hostBaseUrl { // 在远程节点进行处理
				data, err := client.Get(g.name, key)
				if err != nil {
					return nil, err
				}
				return cache.NewByteView(data), nil
			}
		}

		// 先从getter中获取
		data, err := g.getter.Get(key)
		if err != nil {
			return nil, err
		}

		// 获取成功则缓存起来并返回
		byteView := cache.NewByteView(data)
		err = g.safetyCache.Add(key, byteView)
		return byteView, err
	})

	return data.(cache.IByteView), err
}

func (g *Group) Delete(key string) (err error) {
	_, err = g.singleFighting.Do(key+"-DELETE", func() (any, error) {
		// 判断是否开启了分布式节点模式
		if clientsGetter != nil { // 开启了
			client, _ := clientsGetter.GetNodeClient(key)
			if client.GetBaseUrl() != hostBaseUrl { // 在远程节点进行处理
				return nil, client.Delete(g.name, key)
			}
		}
		return nil, g.safetyCache.Delete(key)
	})
	return
}
```

# onEvicted为nil需要特殊处理
- 多加一个onEvicted == nil的特殊处理情况才行，不然下面嵌套掉用的onEvicted是nil就会报错
```go
func NewSafetyCache(maxBytes int64, onEvicted func(key string, value IByteView)) ISafetyCache {
	if onEvicted == nil {
		onEvicted = func(key string, value IByteView) {
			fmt.Printf("execute default onEvicted: the key = %v has been deleted, the value is %v", key, value.String())
		}
	}
	return &safetyCache{
		mu: sync.RWMutex{},
		lruCache: lru.NewLruCache(maxBytes, func(key string, value lru.Value) {
			onEvicted(key, value.(IByteView))
		}),
	}
}
```