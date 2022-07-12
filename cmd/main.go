package main

import (
	"context"
	"sync"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	_map "github.com/atomix/atomix-go-client/pkg/atomix/map"
	"github.com/onosproject/monitor-service/pkg/logger"
	northboundInterface "github.com/onosproject/monitor-service/pkg/northbound"
)

var log = logger.GetLogger()

// Starts some components of the monitor-service
func main() {
	getDistributedMap()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go northboundInterface.Northbound(&waitGroup)

	waitGroup.Wait()
}

func getDistributedMap() {
	ctx := context.Background()

	log.Info("Getting Map")

	myMap, err := atomix.GetMap(ctx, "monitor-config")
	if err != nil {
		log.Errorf("Error from atomix.GetMap: %v\n", err)
		return
	}

	log.Info("Watching myMap for events")

	newChannel := make(chan _map.Event)
	err = myMap.Watch(ctx, newChannel)
	if err != nil {
		log.Errorf("Error getting entry \"Test\" from myMap: %v\n", err)
		return
	}

	go func() {
		select {
		case event := <-newChannel:
			log.Errorf("Event from channel is: %v\n", event)
		}
	}()
}
