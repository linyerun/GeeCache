package cache

// byteView 它永远是只读的, 只要我不大写b, 反射你也别想改我的值
type byteView struct {
	b []byte // 虽然b的空间被很多对象共享，但是它们就是动不了它的值
}

// NewByteView 使用深拷贝，切断b和ByteView的联系
func NewByteView(b []byte) IByteView {
	return byteView{b: CloneBytes(b)}
}

// Len 相等于继承了lru包下的Value接口, ByteView占用字节数
func (v byteView) Len() int {
	return len(v.b)
}

// ByteSlice 获取它的字节数组
func (v byteView) ByteSlice() []byte {
	return CloneBytes(v.b)
}

// String 转string
func (v byteView) String() string {
	return string(v.b) // 转string它是只读的, string再转[]byte会发生深拷贝的
}
