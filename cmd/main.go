package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/atomix/atomix-go-client/pkg/atomix"
	_map "github.com/atomix/atomix-go-client/pkg/atomix/map"
	"github.com/onosproject/monitor-service/pkg/logger"
	northboundInterface "github.com/onosproject/monitor-service/pkg/northbound"
)

// Starts some components of the monitor-service
func main() {
	logger.InitLogging()

	getDistributedMap()

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go northboundInterface.Northbound(&waitGroup)

	waitGroup.Wait()
}

func getDistributedMap() {
	ctx := context.Background()

	fmt.Println("Getting Map")

	myMap, err := atomix.GetMap(ctx, "monitor-config")
	if err != nil {
		// fmt.Printf("Error from atomixClient.GetMap:%+v\n", err)
		fmt.Printf("Error from atomix.GetMap: %v\n", err)
		return
	}

	fmt.Println("Wathing myMap for events")

	newChannel := make(chan _map.Event)
	err = myMap.Watch(ctx, newChannel)
	if err != nil {
		fmt.Printf("Error getting entry \"Test\" from myMap: %v\n", err)
		return
	}

	go func() {
		select {
		case event := <-newChannel:
			fmt.Printf("Event from channel is: %v\n", event)
		}
	}()
}
