package main

import (
	"sync"

	dataProcMgr "github.com/onosproject/monitor-service/pkg/dataProcessingManager"
	northboundInterface "github.com/onosproject/monitor-service/pkg/northbound"

	"github.com/onosproject/monitor-service/pkg/logger"
)

const numberOfModules = 3

// Starts some components of the monitor-service
func main() {
	logger.InitLogging()

	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfModules)

	// streamMgrChannel := make(chan types.StreamMgrChannelMessage)

	go northboundInterface.Northbound(&waitGroup) //, streamMgrChannel)
	go dataProcMgr.DataProcessingManager(&waitGroup)
	// go subscriptionManager.SubscriptionManager(&waitGroup, streamMgrChannel)

	waitGroup.Wait()
}
