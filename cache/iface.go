package cache

type ISafetyCache interface {
	Get(key string) (value IByteView, ok bool)
	Add(key string, value IByteView) (err error) // 新增或删除
	Delete(key string) (err error)
	Size() uint
}

type IByteView interface {
	Len() int
	ByteSlice() []byte
	String() string
}
