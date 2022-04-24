package listeners

import "sync"

type KeyedListeners struct {
	mut       sync.RWMutex
	listeners map[string]*Listeners
}

func (k *KeyedListeners) RegisterListener(key string) <-chan interface{} {
	k.mut.Lock()
	defer k.mut.Unlock()

	listeners, ok := k.listeners[key]
	if !ok {
		listeners = &Listeners{}
		k.listeners[key] = listeners
	}
	return listeners.RegisterListener()
}

func (k *KeyedListeners) UnregisterListener(key string, c <-chan interface{}) {
	k.mut.Lock()
	defer k.mut.Unlock()

	listeners, ok := k.listeners[key]
	if !ok {
		return
	}

	listeners.UnregisterListener(c)
}

type Pair struct {
	key       string
	listeners *Listeners
}

func (k *KeyedListeners) Iterate() <-chan Pair {
	c := make(chan Pair)

	go func() {
		for key, value := range k.listeners {
			c <- Pair{key, value}
		}
		close(c)
	}()

	return c
}

func (k *KeyedListeners) Get(key string) (*Listeners, bool) {
	l, ok := k.listeners[key]
	return l, ok
}
