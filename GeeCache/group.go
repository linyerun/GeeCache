package GeeCache

import (
	"github.com/linyerun/GeeCache/cache"
	"github.com/linyerun/GeeCache/single_fighting"
	"github.com/linyerun/GeeCache/utils"
)

type Group struct {
	name           string
	safetyCache    cache.ISafetyCache
	getter         Getter
	singleFighting single_fighting.ISingleFighting
}

func (g *Group) AddOrUpdate(key string, value cache.IByteView) (err error) {
	// 判断是否开启了分布式节点模式
	if clientsGetter != nil { // 开启了
		client, _ := clientsGetter.GetNodeClient(key)
		if client.GetBaseUrl() != hostBaseUrl { // 在远程节点进行处理
			return client.Post(g.name, key, value.ByteSlice())
		}
	}
	return g.safetyCache.Add(key, value)
}

func (g *Group) Get(key string) (cache.IByteView, error) {
	// 判断缓存中是否有对应的key
	if val, ok := g.safetyCache.Get(key); ok {
		utils.Logger().Println("[GeeCache] Hit key is", key)
		return val, nil
	}

	data, err := g.singleFighting.Do(key, func() (any, error) {
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
	_, err = g.singleFighting.Do(key, func() (any, error) {
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
