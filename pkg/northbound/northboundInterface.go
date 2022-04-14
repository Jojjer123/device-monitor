package northboundInterface

import (
	"sync"
	"time"

	types "github.com/onosproject/monitor-service/pkg/types"
)

func Northbound(waitGroup *sync.WaitGroup, adminChannel chan types.AdminChannelMessage) {
	// fmt.Println("AdminInterface started")
	defer waitGroup.Done()

	// var serverWaitGroup sync.WaitGroup

	// var registerFunction func(chan string, *sync.WaitGroup)

	// select {
	// case x := <-adminChannel:
	// 	{
	// 		registerFunction = x.RegisterFunction
	// 	}
	// }

	adminChannelMessage := <-adminChannel

	go startServer(":11161", adminChannelMessage.ExecuteSetCmd)

	// // Starts the gRPC server which will be the external interface.
	// go startServer(&serverWaitGroup, registerFunction)

	// Wait for the gNMI server to exit before exiting admin interface.
	// serverWaitGroup.Wait()

	for {
		time.Sleep(10 * time.Second)
	}
}
