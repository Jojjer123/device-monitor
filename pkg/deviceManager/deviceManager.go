package deviceManager

import (
	"fmt"
	"sync"

	reqBuilder "github.com/onosproject/device-monitor/pkg/requestBuilder"
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

	deviceMonitorWaitGroup.Wait()
	// fmt.Println("Device manager shutting down...")
}

func executeAdminSetCmd(cmd string, target string, configIndex int) string {
	// fmt.Println(cmd)
	switch cmd {
	case "Create":
		{
			fmt.Println("Creating new device monitor for target: " + target)
			// TODO: Build request, then create device monitor with the request.
			requests := reqBuilder.GetRequest(target, "default", configIndex) // confType options: default & temporary

			fmt.Println(requests)
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

func createDeviceMonitor(numberOfDeviceMonitors *int, managerChannel <-chan string, deviceMonitorWaitGroup *sync.WaitGroup) {
	*numberOfDeviceMonitors += 1
	deviceMonitorWaitGroup.Add(1)

	config := "test-configuration"
	go deviceMonitor(config, numberOfDeviceMonitors, managerChannel, deviceMonitorWaitGroup)
}
