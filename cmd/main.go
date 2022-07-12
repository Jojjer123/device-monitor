package main

import (
	"sync"

	"github.com/onosproject/monitor-service/pkg/logger"
	northboundInterface "github.com/onosproject/monitor-service/pkg/northbound"
)

var log = logger.GetLogger()

// Starts Northbound server of the monitor-service
func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go northboundInterface.Northbound(&waitGroup)

	waitGroup.Wait()
}
