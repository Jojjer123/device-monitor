package deviceManager

import (
	"fmt"
	"sync"
)

const maxNumberOfDeviceMonitors = 10

func DeviceManager(waitGroup *sync.WaitGroup, adminChannel <-chan string) {
	fmt.Println("DeviceManager started")
	defer waitGroup.Done()

	var deviceMonitorWaitGroup sync.WaitGroup

	// Create a map of channels as keys to indexes in the managerChannels slice.
	managerChannelsMap := make(map[chan string]int)

	// TODO: Create a map with IP as keys and channels as values.

	// Create a slice for keeping dynamically created channels.
	var managerChannels []chan string

	// Start device manager with 0 device monitors.
	numberOfDeviceMonitors := 0

	deviceManagerIsActive := true
	for deviceManagerIsActive {
		select {
		case msg := <-adminChannel:
			if msg == "shutdown" {
				fmt.Println("Device manager received shutdown command")
				for i := 0; i < len(managerChannels); i++ {
					managerChannels[i] <- msg
				}
				// deviceManagerIsActive = false
			} else if msg == "create new" {
				fmt.Println("Device manager received create new command")
				if numberOfDeviceMonitors < maxNumberOfDeviceMonitors {
					fmt.Println("Create new device monitor...")

					// The following 3 implemented lines could be replaced if there is only one map with IP and channels.
					channelIndexToUse := len(managerChannels)
					managerChannels = append(managerChannels, make(chan string))

					// Add the newly created channel mapped to its index.
					managerChannelsMap[managerChannels[channelIndexToUse]] = channelIndexToUse

					createDeviceMonitor(&numberOfDeviceMonitors, managerChannels[channelIndexToUse], &deviceMonitorWaitGroup)
				}
			}
		}
	}

	deviceMonitorWaitGroup.Wait()
	// fmt.Println("Device manager shutting down...")
}

func createDeviceMonitor(numberOfDeviceMonitors *int, managerChannel <-chan string, deviceMonitorWaitGroup *sync.WaitGroup) {
	*numberOfDeviceMonitors += 1
	deviceMonitorWaitGroup.Add(1)

	config := "test-configuration"
	go deviceMonitor(config, numberOfDeviceMonitors, managerChannel, deviceMonitorWaitGroup)
}
