package lru

type entry struct {
	key   string
	value Value
}

func newEntry(k string, v Value) *entry {
	return &entry{
		key:   k,
		value: v,
	}
}
