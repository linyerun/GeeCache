package lru

type Value interface {
	Len() int
}

type ICache interface {
	Get(key string) (value Value, ok bool)
	Add(key string, value Value) (err error)
	Delete(key string) (err error)
	Size() uint
}
