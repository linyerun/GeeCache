package cache

import (
	"fmt"
	"testing"
)

func TestSafetyCache(t *testing.T) {
	safetyCache := NewSafetyCache(1<<10, func(key string, value IByteView) {
		fmt.Printf("The Key = {%s} has been deleted!\n", key)
	})

	err := safetyCache.Add("hei", NewByteView([]byte("hei, 666")))
	if err != nil {
		t.Error(err)
	}

	value, ok := safetyCache.Get("hei")
	if ok {
		fmt.Println("key = hei, value =\"", value.String(), "\"")
	}

	size := safetyCache.Size()
	fmt.Println("size:", size)

	err = safetyCache.Delete("hei")
	if err != nil {
		t.Error(err)
	}
}

func TestSafetyCache_Get(t *testing.T) {

}

func TestSafetyCache_Add(t *testing.T) {

}

func TestSafetyCache_Delete(t *testing.T) {

}

func TestSafetyCache_Size(t *testing.T) {

}
