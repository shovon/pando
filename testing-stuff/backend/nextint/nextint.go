package nextint

import "sync"

var value int64 = 0

var m sync.Mutex

func NextInt() int64 {
	m.Lock()
	defer m.Unlock()
	value++
	return value
}
