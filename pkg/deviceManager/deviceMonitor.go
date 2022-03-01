package deviceManager

import (
	"fmt"
	"sync"
	"time"
)

func deviceMonitor(config string, numberOfDeviceMonitors *int, managerChannel <-chan string, deviceMonitorWaitGroup *sync.WaitGroup) {
	defer deviceMonitorWaitGroup.Done()
	var localWaitGroup sync.WaitGroup
	localWaitGroup.Add(1)

	// TODO: Need to have a way to communicate with specific goroutines for specific counters

	// Temporary interval, should be specified in the config?
	interval := 2

	// // Create a map of counter names as keys to channel references in the counterChannels slice.
	// counterChannelsMap := make(map[string]chan string)

	// // Create a slice for keeping dynamically created channels.
	// var counterChannels []chan string

	// This should happen for every counter we want to have
	deviceMonitorChannel := make(chan string)
	go newCounter(config, interval, &localWaitGroup, deviceMonitorChannel)
	fmt.Println("Created new device monitor")

	// Loops forever and reads the channel which admin interface controls, if the channel has any new data, read it and react accordingly.
	deviceMonitorIsActive := true
	for deviceMonitorIsActive {
		select {
		case x := <-managerChannel:
			if x == "shutdown" { // shut down (all) device monitors
				fmt.Println("Received shutdown command on channel now...")
				deviceMonitorIsActive = false
				deviceMonitorChannel <- x
				*numberOfDeviceMonitors -= 1
			} else if x == "update" {
				fmt.Println("Received update command on channel now...")

			}
		}
	}

	localWaitGroup.Wait()
	fmt.Println("Shutting down device monitor now...")
}

func sendToChannel(msg string) {
	// TODO: Check which channel should get the msg and send it on that channel.
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
			fmt.Println("Ticker triggered")
		}
	}

	fmt.Println("Exits counter now")
}
