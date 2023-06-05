package single_fighting

import "sync"

type singleFighting struct {
	resMapRwMutex sync.RWMutex
	resMap        map[string]*callResult
}

func NewSingleFighting() ISingleFighting {
	return &singleFighting{
		resMap: make(map[string]*callResult),
	}
}

func (s *singleFighting) Do(uniqueMark string, fn func() (any, error)) (any, error) {
	//读
	s.resMapRwMutex.RLock()
	res, ok := s.resMap[uniqueMark]
	s.resMapRwMutex.RUnlock()

	if ok { // 已经请求在调用了, 我们等着它共享数据即可
		res.wg.Wait()
		return res.value, res.err
	}

	// 这是首次调用
	res = new(callResult)
	res.wg.Add(1)

	// 写
	s.resMapRwMutex.Lock()
	s.resMap[uniqueMark] = res
	s.resMapRwMutex.Unlock()

	res.value, res.err = fn() // 调用
	res.wg.Done()

	// 删
	s.resMapRwMutex.Lock()
	delete(s.resMap, uniqueMark)
	s.resMapRwMutex.Unlock()

	return res.value, res.err
}
