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
	managerChannel := make(chan string)
	numberOfDeviceMonitors := 0

	deviceManagerIsActive := true
	for deviceManagerIsActive {
		select {
		case msg := <-adminChannel:
			if msg == "shutdown" {
				fmt.Println("Device manager received shutdown command")
				// deviceManagerIsActive = false
				for i := 0; i < numberOfDeviceMonitors; i++ {
					managerChannel <- msg
				}
			} else if msg == "create new" {
				fmt.Println("Device manager received create new command")
				if numberOfDeviceMonitors < maxNumberOfDeviceMonitors {
					fmt.Println("Create new device monitor...")
					createDeviceMonitor(&numberOfDeviceMonitors, managerChannel, &deviceMonitorWaitGroup)
				}
			}
		}
	}

	deviceMonitorWaitGroup.Wait()
	fmt.Println("Device manager shutting down...")
}

func createDeviceMonitor(numberOfDeviceMonitors *int, managerChannel <-chan string, deviceMonitorWaitGroup *sync.WaitGroup) {
	*numberOfDeviceMonitors += 1
	deviceMonitorWaitGroup.Add(1)

	config := "test-configuration"
	go deviceMonitor(config, numberOfDeviceMonitors, managerChannel, deviceMonitorWaitGroup)
}
