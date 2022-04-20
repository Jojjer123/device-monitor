package northboundInterface

import (
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/types"
)

func Northbound(waitGroup *sync.WaitGroup, adminChannel chan types.ConfigAdminChannelMessage, streamMgrChannel chan types.StreamMgrChannelMessage) {
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
	streamMgrChannelMessage := <-streamMgrChannel

	go startServer(":11161", adminChannelMessage.ExecuteSetCmd, streamMgrChannelMessage.ManageCmd)

	// // Starts the gRPC server which will be the external interface.
	// go startServer(&serverWaitGroup, registerFunction)

	// Wait for the gNMI server to exit before exiting admin interface.
	// serverWaitGroup.Wait()

	for {
		time.Sleep(10 * time.Second)
	}
}
