package northboundInterface

import (
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/types"
)

func Northbound(waitGroup *sync.WaitGroup, adminChannel chan types.ConfigAdminChannelMessage, streamMgrChannel chan types.StreamMgrChannelMessage) {
	defer waitGroup.Done()

	adminChannelMessage := <-adminChannel
	streamMgrChannelMessage := <-streamMgrChannel

	// Starts a gRPC server.
	go startServer(false, ":11161", adminChannelMessage.ExecuteSetCmd, streamMgrChannelMessage.ManageCmd)
	go startServer(true, ":10161", adminChannelMessage.ExecuteSetCmd, streamMgrChannelMessage.ManageCmd)

	// Remove???
	for {
		time.Sleep(10 * time.Second)
	}
}
