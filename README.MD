# 快速入门
```go
package main

import (
	"flag"
	"fmt"
	"github.com/linyerun/GeeCache/GeeCache"
)

func main() {
	var port uint
	flag.UintVar(&port, "port", 0, "GeeCache server port")
	flag.Parse()

	// 引入单体项目作为一个Cache模块使用
	g := GeeCache.NewGroupByGetterFunc("students", 1<<20, func(key string) ([]byte, error) {
		return []byte("can not find key=" + key), nil
	}, nil)

	group, ok := GeeCache.GetGroup("student")
	if ok {
		fmt.Println(g == group)
	}

	GeeCache.NewGroupByGetterFunc("dogs", 1<<20, func(key string) ([]byte, error) {
		return []byte("can not find key=" + key), nil
	}, nil)

	GeeCache.NewGroupByGetterFunc("cats", 1<<20, func(key string) ([]byte, error) {
		return []byte("can not find key=" + key), nil
	}, nil)

	// 开启分布式
	httpClientGetter := GeeCache.NewHttpClientGetter(20, nil)
	httpClientGetter.Set("127.0.0.1:6666", "127.0.0.1:9999", "127.0.0.1:8888", "127.0.0.1:7777")

	// 对外HTTP服务提供
	server := GeeCache.NewHttpServer("127.0.0.1", port, "/")
	if err := server.Run(); err != nil {
		panic(err)
	}

}
```

```shell
#! /bin/bash

go build -o server main.go

./server -port 6666 &
./server -port 7777 &
./server -port 8888 &
./server -port 9999 &

rm -f server
```

# utils包
- 目前只整合了logrus第三方包进来，作用是进行日志输出

# lru包介绍
- 实现了LruCache, key为string, value是实现了Value接口的对象
- 支持操作的接口:
  - Get: 根据key获取Value 
  - Add: 如果key不存在则新增(key,value); 存在则更新key对应的值为value参数对应的值
  - Delete: 删除lruCache中key对应的KeyValue对
- 部分属性介绍:
  - onEvicted: lruCache的KV对被删除是执行的函数
  - maxCacheBytes: lruCache最大占用的字节数(加进的Key和Value占用的总字节)
  - curCacheBytes: lruCache当前占用的字节数

# cache包介绍
- 因为lru包实习的lruCache是并发不安全的，所以cache在lru包的基础上对内容进行了加锁封装，得到safetyCache

# consistent_hash包介绍
- 实现一致性哈希算法，把我们的单机缓存扩展成分布式缓存。
- 作用：随意一个请求打到分布式缓存任意节点上，根据需要寻找的key和consistent_hash包的信息找到需要去哪个节点找
- 属性介绍
  - hashFunc: 一致性哈希算法, []byte => uint32
  - replicas: 虚拟节点 = 真实节点*replicas, 多做一些虚拟节点, 弥补真实节点数量少是在哈希环上分布不均的问题
  - hashMap: 一个Map, key是虚拟节点对应的哈希值, Value是字符串, 真实节点信息
- 注意
  - 如果不同真实节点的虚拟节点对应的hash值一致, 由于map的key具有唯一性, 那么这个虚拟节点只会归属最后来的真实
  - 这也告诫我们在使用这个包新增节点时，一定要保证所有分布式节点新增顺序是一致的，不然导致的后果是分布式节点信息不一致

# single_fighting包介绍
- 作用: 多个幂等的请求达到分布式集群上, 只执行首个请求, 然后首个请求前执行完前的这些幂等请求共享首个请求的结果
- 好处: 避免幂等请求被冗余执行

# protobuf包的介绍
- 使用protobuf协议序列化发送的HTTP数据
- 主要针对两个方面的: Request和Response
```protobuf
syntax = "proto3";

package protobuf;
option go_package="../protobuf";

message Request {
  string group = 1; // groupName
  string key = 2;   // key
}

message Response {
  bytes value = 1;  // 字节数组
}

service GroupCache {
  rpc Get(Request) returns (Response);
}
```

# GeeCache包介绍
- 这个包是框架最主要的部分, 我们一个一个文件地来进行分析

## for_user.go文件的内容
- 这是提供给框架的使用者使用的构造函数
- 开发者通过这些构造函数可以使用一下功能:
  - 开启缓存分组
  - 开启HTTP服务对外暴露节点, 使得外部应用可以使用该节点的缓存分组
  - 开启分布式节点服务
## global.go文件的内容
- 存放GeeCache包被的私有全局遍历
- clientsGetter: 开启分布式存储后，该变量存储的是所有节点发送请求的客户端
- groupsRwMutex, groups: 用户并发安全的获取缓存分组元素来执行一些操作
- hostBaseUrl: 当前节点ip+port的组合而形成的字符串
- binaryContentType: 常量，指定响应报文响应体的内容形式
## group.go文件内容
```go
type Group struct {
	name           string
	safetyCache    cache.ISafetyCache
	getter         Getter
	singleFighting single_fighting.ISingleFighting
}
```
- 缓存组的属性: 名称、集成safetyCache、getter在缓存获取不到某个KV对时执行(配置时为空会报错)、singleFighting用于避免缓存击穿
- 它的功能: 提供了对缓存组的增删改查功能，一个key过来，如果开启了分布式缓存会使用全局变量clientsGetter获取一个httpClient来发启远程请求，当然前提是当前节点不具备处理该key的资格

## http_client文件内容
```go
type httpClient struct {
	baseUrl string
}
func jointServerURL(baseUrl, group, key string) string {
    return fmt.Sprintf(
      "http://%v%v%v/%v",
      baseUrl,
      globalBasePath,
      url.QueryEscape(group), //QueryEscape函数对s进行转码使之可以安全的用在URL查询里。有些特殊字符不能用可转的
      url.QueryEscape(key),
    )
}
```
- httpClient提供了GET、POST、DELETE三种请求方法，因为GET请求获取到的数据使用了Protobuf协议，所以会进行Unmarshal操作然后返回结果
- 遗憾的是我们并没有使用我们生成在protobuf包下的Request对key和groupName进行二进制化。

## http_clients_getter
```go
type httpClientsGetter struct {
	nodes       consistent_hash.IConsistentHash
	httpClients map[string]INodeClient
}
```
- 它整合了consistent_hash的一致性哈希算法
- 作用: 根据key获取可以处理它的客户端httpClient

## http_server
- 用于处理当前节点的http请求，提供了GET、POST、DELETE三种处理请求的逻辑，最后它使用的是Group来获取结果响应请求。