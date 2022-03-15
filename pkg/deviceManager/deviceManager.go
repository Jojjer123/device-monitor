package deviceManager

import (
	"fmt"
	"sync"

	types "github.com/onosproject/device-monitor/pkg/types"
)

// const maxNumberOfDeviceMonitors = 10

func DeviceManager(waitGroup *sync.WaitGroup, adminChannel chan types.AdminChannelMessage) {
	fmt.Println("DeviceManager started")
	defer waitGroup.Done()

	var deviceMonitorWaitGroup sync.WaitGroup

	var adminMessage types.AdminChannelMessage
	adminMessage.RegisterFunction = registerServerChannel
	adminMessage.ExecuteSetCmd = executeAdminSetCmd

	adminChannel <- adminMessage

	fmt.Println(<-adminChannel)
	// deviceManagerIsActive := true
	// for deviceManagerIsActive {

	// 	select {
	// 	case msg := <-adminChannel:
	// 		{
	// 			// deviceManagerIsActive = false
	// 			fmt.Println(msg)
	// 		}
	// 	}
	// }

	if false {
		go createDeviceMonitor(nil, nil, nil)
	}

	// // Create a map of channels as keys to indexes in the managerChannels slice.
	// managerChannelsMap := make(map[chan string]int)

	// // TODO: Create a map with IP as keys and channels as values.

	// // Create a slice for keeping dynamically created channels.
	// var managerChannels []chan string

	// // Start device manager with 0 device monitors.
	// numberOfDeviceMonitors := 0

	// TODO: Rework this for-loop, there should be a function for registering new channels & this loop
	//		 will be used in the new function.
	// deviceManagerIsActive := true
	// for deviceManagerIsActive {
	// 	select {
	// 	case msg := <-adminChannel:
	// 		if msg == "shutdown" {
	// 			fmt.Println("Device manager received shutdown command")
	// 			for i := 0; i < len(managerChannels); i++ {
	// 				managerChannels[i] <- msg
	// 			}
	// 			// deviceManagerIsActive = false
	// 		} else if msg == "create new" {
	// 			fmt.Println("Device manager received create new command")
	// 			if numberOfDeviceMonitors < maxNumberOfDeviceMonitors {
	// 				fmt.Println("Create new device monitor...")

	// 				// The following 3 implemented lines could be replaced if there is only one map with IP and channels.
	// 				channelIndexToUse := len(managerChannels)
	// 				managerChannels = append(managerChannels, make(chan string))

	// 				// Add the newly created channel mapped to its index.
	// 				managerChannelsMap[managerChannels[channelIndexToUse]] = channelIndexToUse

	// 				createDeviceMonitor(&numberOfDeviceMonitors, managerChannels[channelIndexToUse], &deviceMonitorWaitGroup)
	// 			}
	// 		}
	// 	}
	// }

	deviceMonitorWaitGroup.Wait()
	// fmt.Println("Device manager shutting down...")
}

func executeAdminSetCmd(cmd string) string {
	fmt.Println(cmd)
	// switch cmdType {
	// case "":

	// }
	return "Successfully executed command sent"
}

// TODO: REWORK THIS, could potentially be used as a basic function, call it with msg as param and return string.
// Executes the request/message and returns the response on the provided channel.
func registerServerChannel(serverChannel chan string, channelWaitGroup *sync.WaitGroup) {
	defer channelWaitGroup.Done()
	fmt.Println(<-serverChannel)
	// select {
	// case x := <-serverChannel:
	// 	{
	// 		fmt.Println(x)
	// 	}
	// }

	// TODO: Create, update or delete device monitor and get a response from it.
	response := "Success"
	serverChannel <- response
}

func createDeviceMonitor(numberOfDeviceMonitors *int, managerChannel <-chan string, deviceMonitorWaitGroup *sync.WaitGroup) {
	*numberOfDeviceMonitors += 1
	deviceMonitorWaitGroup.Add(1)

	config := "test-configuration"
	go deviceMonitor(config, numberOfDeviceMonitors, managerChannel, deviceMonitorWaitGroup)
}
