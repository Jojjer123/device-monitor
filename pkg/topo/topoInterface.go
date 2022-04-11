package topo

import (
	"sync"
)

func TopoInterface(waitGroup *sync.WaitGroup) {
	// fmt.Println("Topo interface started")
	defer waitGroup.Done()
}
