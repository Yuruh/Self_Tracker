package spotify

import (
	"sync"
	"sync/atomic"
)

var mu sync.Mutex
var initialized uint32 = 0
var instance *Connector

type Connector struct {}


func GetConnector() *Connector {
	if atomic.LoadUint32(&initialized) == 1 {
		return instance
	}
	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		instance = &Connector{
		}
		atomic.StoreUint32(&initialized, 1)
	}

	return instance
}

func (spotify *Connector) Ping() {
	println("PONG")
}
