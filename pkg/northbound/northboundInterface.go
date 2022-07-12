package northboundInterface

import (
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/logger"
	// "github.com/onosproject/monitor-service/pkg/types"
)

var log = logger.GetLogger()

func Northbound(waitGroup *sync.WaitGroup) { //, streamMgrChannel chan types.StreamMgrChannelMessage) {
	defer waitGroup.Done()

	// streamMgrChannelMessage := <-streamMgrChannel

	// Starts a gRPC server.
	go startServer(false, ":11161") //, streamMgrChannelMessage.ManageCmd)
	go startServer(true, ":10161")  //, streamMgrChannelMessage.ManageCmd)

	// Remove???
	for {
		time.Sleep(10 * time.Second)
	}
}
