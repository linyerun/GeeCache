# demo01: 实现LruCache
## 目标
- 实现一个基于LRU淘汰机制的Cache
## 注意点
- logrus关于打印日志的方法都不需要加 `\n`, 它自动帮我们加了的。
## 测试
- 参考文章：https://geektutu.com/post/quick-go-test.html

# demo02: 实现缓存分组
## 并发
- 因为需要对Cache进行并发访问，所以我们要对Cache进行封装, 封装成SafetyCache
## byteView和IByteView
- 我们希望byteView的字节数组不会被用户篡改，我们返回它的接口就不怕了，byteView的方法也保护它的数据不会泄露
## 分组
- 全局维护一个groups的map，和一个读写锁保护groups的修改和访问
- 提供了两个创建Group的函数和一个获取指定Group的函数

# demo03: 实现HTTP服务端，对外提供缓存服务
## IServer
- 只要实现这个接口，我们可以基于TCP通信、HTTP通信都可以的

# demo04: 实现一致性Hash算法
## 目标
- 通过该算法来确定我们的请求应该达到哪个节点上。

# demo05: 实现分布式节点
## 目标
- 借助上一天封装好的consistentHash对象帮我们找到所需的服务端。
## 注意点
- 在单机情况下，分组只存在于一台服务器上，当我们使用分布式节点之后，每个节点都有所有的分组，但是分组中的key分布在不同的节点上
- 分组必须在启动分布式服务器之前初始化好，后续不能再初始化了。
## 实现
现在统一让更新或者新增使用post请求即可
```go
func (h *httpClient) Put(group string, key string, value []byte) error {
	serverURL := jointServerURL(h.baseUrl, group, key)

	// 创建HttpRequest
	request, err := http.NewRequest(http.MethodPut, serverURL, bytes.NewBuffer(value))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", binaryContentType)

	// 发送请求
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			utils.Logger().Errorln(err.Error())
			return err
		}
		return errors.New(string(data))
	}

	return nil
}
```

# demo06
- 实现了一个singleFighting结构体，它的作用是给每个group对象都配备上，当group对象想让自己的某个方法在一段时间内收到了很多幂等请求，那就只处理第一个请求，让其他请求共享第一个请求的结果。
- 好处：可以避免不必要的HTTP请求服务或者寻找数据库的操作。
- 在single_fighting包实现了，并整合到GeeCache的group结构体中。

# demo07
- 使用Protobuf压缩数据