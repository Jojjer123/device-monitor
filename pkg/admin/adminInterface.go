package adminInterface

import (
	"fmt"
	"sync"
	"time"
)

func AdminInterface(waitGroup *sync.WaitGroup, adminChannel chan<- string) {
	fmt.Println("AdminInterface started")
	defer waitGroup.Done()

	time.Sleep(2 * time.Second)

	// Communicate with device manager over deviceChannel.
	adminChannel <- "create new"
	adminChannel <- "create new"

	time.Sleep(10 * time.Second)
	fmt.Println("Send shutdown command over channel now...")
	adminChannel <- "shutdown"
}
