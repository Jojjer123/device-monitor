package deviceManager

import (
	"fmt"
	"sync"
	"time"
)

func deviceMonitor(config string, numberOfDeviceMonitors *int, deviceChannel <-chan string) {
	var localWaitGroup sync.WaitGroup
	localWaitGroup.Add(1)

	// TODO: Need to have a way to communicate with specific goroutines for specific counters
	// TODO: Send message to goroutines dedicated for each counter, to change already existing device monitor.

	// For each interval, create a new goroutine, which repeatedly collects counters.
	interval := 2
	deviceMonitorChannel := make(chan string)
	go newCounter(config, interval, &localWaitGroup, deviceMonitorChannel)
	fmt.Println("Created new device monitor")

	// Loops forever and reads the channel which admin interface controls, if the channel has any new data, read it and react accordingly.
	deviceMonitorIsActive := true
	for deviceMonitorIsActive {
		select {
		case x := <-deviceChannel:
			if x == "shutdown" {
				// fmt.Println("Received shutdown command on channel now...")
				deviceMonitorIsActive = false
				deviceMonitorChannel <- x
				*numberOfDeviceMonitors -= 1
			}
		}
	}
	// fmt.Println("Shutting down device monitor now...")
}

func newCounter(config string, interval int, localWaitGroup *sync.WaitGroup, deviceMonitorChannel <-chan string) {
	defer localWaitGroup.Done()
	// Start a ticker which will trigger repeatedly after (interval) seconds.
	intervalTicker := time.NewTicker(time.Duration(interval*1000) * time.Millisecond)

	counterIsActive := true
	for counterIsActive == true {
		select {
		case <-deviceMonitorChannel:
			counterIsActive = false
		case <-intervalTicker.C:
			// TODO: Get the counters for the given interval here and send them to the data processing part.
			// fmt.Println("Ticker triggered")
		}
	}

	// fmt.Println("Exits counter now")
}
