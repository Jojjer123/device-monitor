package northboundInterface

import (
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/types"
)

func Northbound(waitGroup *sync.WaitGroup, adminChannel chan types.ConfigAdminChannelMessage, streamMgrChannel chan types.StreamMgrChannelMessage) {
	// fmt.Println("AdminInterface started")
	defer waitGroup.Done()

	adminChannelMessage := <-adminChannel
	streamMgrChannelMessage := <-streamMgrChannel

	// Starts a gRPC server.
	go startServer(":11161", adminChannelMessage.ExecuteSetCmd, streamMgrChannelMessage.ManageCmd)

	// TODO: Create another gRPC server on secure port (10161).

	// Remove???
	for {
		time.Sleep(10 * time.Second)
	}
}
