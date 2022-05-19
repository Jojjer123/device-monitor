package main

import (
	"sync"

	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/northbound"
)

// Starts some components of the monitor-service
func main() {
	logger.InitLogging()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go northboundInterface.Northbound(&waitGroup)

	waitGroup.Wait()
}
