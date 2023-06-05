package single_fighting

import "sync"

type callResult struct {
	wg    sync.WaitGroup
	value any
	err   error
}
