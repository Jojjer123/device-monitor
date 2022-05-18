package main

import (
	"sync"

	dataProcMgr "github.com/onosproject/monitor-service/pkg/dataProcessingManager"
	north "github.com/onosproject/monitor-service/pkg/northbound"
	streamMgr "github.com/onosproject/monitor-service/pkg/streamManager"

	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/types"
)

const numberOfModules = 3

// Starts some components of the monitor-service
func main() {
	logger.InitLogging()

	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfModules)

	streamMgrChannel := make(chan types.StreamMgrChannelMessage)

	go north.Northbound(&waitGroup, streamMgrChannel)
	go dataProcMgr.DataProcessingManager(&waitGroup)
	go streamMgr.StreamManager(&waitGroup, streamMgrChannel)

	waitGroup.Wait()
}
