package northboundInterface

import (
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/logger"
)

var log = logger.GetLogger()

func Northbound(waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Starts two gNMI servers
	go startServer(false, ":11161")
	go startServer(true, ":10161")

	// TODO: Replace with proper methods, such as a sync.WaitGroup
	for {
		time.Sleep(10 * time.Second)
	}
}
