package lru

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/linyerun/GeeCache/utils"
)

var logger = utils.Logger()

type lruCache struct {
	maxCacheBytes int64 //最大可以缓存多少字节
	curCacheBytes int64 //当前缓存的多少字节
	lruList       *list.List
	lruMap        map[string]*list.Element
	onEvicted     func(key string, value Value) //可选的，在条目被清除时执行。evicted：驱逐
}

func NewLruCache(maxCacheBytes int64, onEvicted func(key string, value Value)) ICache {
	return &lruCache{
		maxCacheBytes: maxCacheBytes,
		curCacheBytes: 0,
		lruList:       list.New(),
		lruMap:        map[string]*list.Element{},
		onEvicted:     onEvicted,
	}
}

// Get 获取元素
func (c *lruCache) Get(key string) (value Value, ok bool) {
	if elem, ok := c.lruMap[key]; ok {
		c.lruList.MoveToFront(elem)
		logger.Printf("Get the value of the key = %v successfully.", key)
		return elem.Value.(*entry).value, true
	}
	return
}

// removeOldest 移除最近最少被使用的元素
func (c *lruCache) removeOldest() {
	lruList := c.lruList

	//获取最后一个元素
	revElem := lruList.Back()
	if revElem == nil { //无元素则之间返回
		return
	}

	//删除
	lruList.Remove(revElem)
	e := revElem.Value.(*entry)
	delete(c.lruMap, e.key)
	c.curCacheBytes -= int64(len(e.key) + e.value.Len())
	logger.Infof("remove [key=%s,value=%v] successfully,current bytes is %d", e.key, e.value, c.curCacheBytes)

	//驱逐元素之后可以执行的操作
	if c.onEvicted != nil {
		logger.Printf("begin to call onEvicted![key=%s,value=%s]", e.key, e.value)
		c.onEvicted(e.key, e.value)
	}
}

// Add 缓存添加
func (c *lruCache) Add(key string, value Value) (err error) {
	// 判断新加进来的kv是否大于maxBytes
	if c.maxCacheBytes < int64(value.Len()+len(key)) {
		errMsg := fmt.Sprintf("Your value of your key = %v is too large to cache.", key)
		logger.Errorf(errMsg)
		return errors.New(errMsg)
	}

	if elem, ok := c.lruMap[key]; ok { //更新操作
		c.lruList.MoveToFront(elem)                            //移动到头
		kv := elem.Value.(*entry)                              //类型断言转kv
		c.curCacheBytes += int64(value.Len() - kv.value.Len()) //更新curBytes
		kv.value = value                                       //更新值
		logger.Infof("Update [key=%s,value=%v] successfully,current bytes is %d", key, value, c.curCacheBytes)
	} else { //添加操作
		c.lruMap[key] = c.lruList.PushFront(newEntry(key, value)) //加入链表中、加入map中
		c.curCacheBytes += int64(len(key) + value.Len())          //更改curBytes值
		logger.Infof("Add [key=%s,value=%v] successfully,current bytes is %d", key, value, c.curCacheBytes)
	}

	//清除缓存，知道maxBytes >= curCacheBytes
	for c.maxCacheBytes < c.curCacheBytes { // maxCacheBytes >= curCacheBytes 结束条件
		c.removeOldest()
	}
	return
}

// Size 当前缓存个数
func (c *lruCache) Size() uint {
	return uint(c.lruList.Len())
}

func (c *lruCache) Delete(key string) (err error) {
	// 判断是否存在
	elem, ok := c.lruMap[key]
	e := elem.Value.(*entry)
	if !ok {
		return errors.New("this key cannot be found in the cache system")
	}

	// 在lruMap中删除
	delete(c.lruMap, key)

	// 在lruList中删除
	c.lruList.Remove(elem)

	// 更新curBytes值
	c.curCacheBytes -= int64(len(key) + e.value.Len())

	//驱逐元素之后可以执行的操作
	if c.onEvicted != nil {
		logger.Printf("begin to call onEvicted![key=%s,value=%s]", e.key, e.value)
		c.onEvicted(e.key, e.value)
	}
	return
}
