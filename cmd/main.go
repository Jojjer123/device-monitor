package main

import (
	"sync"

	northboundInterface "github.com/onosproject/monitor-service/pkg/northbound"

	"github.com/onosproject/monitor-service/pkg/logger"
)

// Starts some components of the monitor-service
func main() {
	logger.InitLogging()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go northboundInterface.Northbound(&waitGroup)
	// go dataProcMgr.DataProcessingManager(&waitGroup)

	waitGroup.Wait()
}
