package dataProcessing

import (
	"sync"
)

func DataProcessingManager(waitGroup *sync.WaitGroup) {
	// fmt.Println("DataProcessing started")
	defer waitGroup.Done()
}
