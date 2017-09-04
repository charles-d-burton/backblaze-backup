package datastores

import (
	"sync"
)

type accumulator struct {
	sync.Mutex
	Files map[string]bool
}

var accumulatorInstance *accumulator
var once sync.Once

func GetAccumulator() *accumulator {
	once.Do(func() {
		accumulatorInstance = &accumulator{}
	})
	return accumulatorInstance
}
