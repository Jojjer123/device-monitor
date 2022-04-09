package deviceManager

import (
	"fmt"
	"sync"

	reqBuilder "github.com/onosproject/device-monitor/pkg/requestBuilder"
	"github.com/onosproject/device-monitor/pkg/types"
)

// const maxNumberOfDeviceMonitors = 10

var deviceMonitorStore []types.DeviceMonitor

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

	// if false {
	// 	go createDeviceMonitor(nil, nil, nil)
	// }

	deviceMonitorWaitGroup.Wait()
	// fmt.Println("Device manager shutting down...")
}

func executeAdminSetCmd(cmd string, target string, configIndex int) string {
	// fmt.Println(cmd)
	switch cmd {
	case "Create":
		{
			// Get slice of the different paths with their intervals and the appropriate
			// adapter if one is necessary
			requests, adapter, target := reqBuilder.GetConfig(target, configIndex)

			fmt.Println(requests)
			fmt.Println(adapter)
			fmt.Println(target)

			// TODO: Create and register a device-monitor in a table?
			createDeviceMonitor(requests, adapter, target)
		}
	case "Update": // Should not be implemented before discussed design (how it should update configs).
		{
			fmt.Println("Updating device monitor with target: " + target)
		}
	default:
		{
			fmt.Println("Could not find command: " + cmd)
			return "Command not found!"
		}
	}

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

func createDeviceMonitor(requests []types.Request, adapter types.Adapter, target string) {
	managerChannel := make(chan string)

	deviceMonitorStore = append(deviceMonitorStore, types.DeviceMonitor{
		Target:         target,
		Adapter:        adapter,
		ManagerChannel: managerChannel,
	})

	go deviceMonitor(target, adapter, requests, managerChannel)

	// *numberOfDeviceMonitors += 1
	// deviceMonitorWaitGroup.Add(1)

	// config := "test-configuration"
	// go deviceMonitor(config, numberOfDeviceMonitors, managerChannel, deviceMonitorWaitGroup)
}
