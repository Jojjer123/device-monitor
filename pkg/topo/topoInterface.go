package topo

import (
	"fmt"
	"sync"
)

func TopoInterface(waitGroup *sync.WaitGroup) {
	fmt.Println("Topo interface started")
	defer waitGroup.Done()
}
