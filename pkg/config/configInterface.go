package configInterface

import (
	"fmt"
	"sync"
)

func ConfigInterface(waitGroup *sync.WaitGroup) {
	fmt.Println("ConfigInterface started")
	defer waitGroup.Done()
}
