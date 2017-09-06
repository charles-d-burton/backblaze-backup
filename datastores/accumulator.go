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

//GetAccumulator ...Return a reference to the accumulator
func GetAccumulator() *accumulator {
	once.Do(func() {
		accumulatorInstance = &accumulator{}
		accumulatorInstance.Files = make(map[string]bool)
	})
	return accumulatorInstance
}
