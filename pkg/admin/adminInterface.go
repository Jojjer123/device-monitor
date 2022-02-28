package adminInterface

import (
	"fmt"
	"sync"
	"time"
)

func AdminInterface(waitGroup *sync.WaitGroup, deviceChannel chan<- string) {
	fmt.Println("AdminInterface started")
	defer waitGroup.Done()

	// Communicate with device manager over deviceChannel.
	time.Sleep(10 * time.Second)
	// fmt.Println("Send shutdown command over channel now...")
	deviceChannel <- "shutdown"
}
