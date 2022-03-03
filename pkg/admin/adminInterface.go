package adminInterface

import (
	"fmt"
	"sync"

	types "github.com/onosproject/device-monitor/pkg/types"
)

func AdminInterface(waitGroup *sync.WaitGroup, adminChannel chan types.AdminChannelMessage) {
	fmt.Println("AdminInterface started")
	defer waitGroup.Done()

	// TODO: Instead of having a channel to gRPC server, get the function for registering channels in device manager,
	//		 from the admin channel, and use that function directly in the gRPC server.

	var serverWaitGroup sync.WaitGroup

	// serverChannel := make(chan string)

	var registerFunction func(chan string, *sync.WaitGroup)

	select {
	case x := <-adminChannel:
		{
			registerFunction = x.RegisterFunction
		}
	}

	serverWaitGroup.Add(1)

	// Starts the gRPC server which will be the external interface.
	go startServer(&serverWaitGroup, registerFunction) //, serverChannel)

	// // Loops forever and reads the channel which gRPC server uses to register new channels with device manager.
	// adminInterfaceIsActive := true
	// for adminInterfaceIsActive {
	// 	select {
	// 	case <-serverChannel:
	// 		// adminInterfaceIsActive = false

	// 		// TODO: Return a bidirectional channel to device manager.
	// 	}
	// }

	/* (WAS USED FOR TESTING)
	* time.Sleep(10 * time.Second)
	*
	* // Communicate with device manager over deviceChannel.
	* adminChannel <- "create new"
	* adminChannel <- "create new"
	*
	* time.Sleep(10 * time.Second)
	* fmt.Println("Send shutdown command over channel now...")
	* adminChannel <- "shutdown"
	 */

	// Wait for the gRPC server to exit before exiting admin interface.
	serverWaitGroup.Wait()
}
