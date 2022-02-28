package deviceManager

import (
	"fmt"
	"sync"
)

const maxNumberOfDeviceMonitors = 10

func DeviceManager(waitGroup *sync.WaitGroup, deviceChannel <-chan string) {
	fmt.Println("DeviceManager started")
	defer waitGroup.Done()

	numberOfDeviceMonitors := 0

	i := 0
	for {
		if i == 0 && numberOfDeviceMonitors < maxNumberOfDeviceMonitors {
			createDeviceMonitor(&numberOfDeviceMonitors, deviceChannel)
		}

		i++
	}
}

func createDeviceMonitor(numberOfDeviceMonitors *int, deviceChannel <-chan string) {
	*numberOfDeviceMonitors += 1

	config := "test-configuration"
	go deviceMonitor(config, numberOfDeviceMonitors, deviceChannel)
}
