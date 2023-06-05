package single_fighting

import (
	"fmt"
	"testing"
	"time"
)

func TestSingleFighting_Do(t *testing.T) {
	singleFighting := NewSingleFighting()
	for i := 0; i < 10; i++ {
		go func(i int) {
			res, _ := singleFighting.Do("666", func() (any, error) {
				time.Sleep(time.Second * 2)
				fmt.Println("我是首次进来的,i =", i)
				return 0, nil
			})
			fmt.Println("i =", i, "res = ", res)
		}(i)
	}
	time.Sleep(time.Second * 5)
}
