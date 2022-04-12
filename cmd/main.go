package main

import (
	"sync"

	confMgr "github.com/onosproject/device-monitor/pkg/configManager"
	dataProcMgr "github.com/onosproject/device-monitor/pkg/dataProcessingManager"
	north "github.com/onosproject/device-monitor/pkg/northbound"
	reqBuilder "github.com/onosproject/device-monitor/pkg/requestBuilder"
	storage "github.com/onosproject/device-monitor/pkg/storage"

	types "github.com/onosproject/device-monitor/pkg/types"
)

const numberOfComponents = 6

// Starts the main components of the device-monitor
func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfComponents)

	// WARNING potential problem: buffered vs unbuffered channels block in different stages of the communication.
	adminChannel := make(chan types.AdminChannelMessage)

	go storage.ConfigInterface(&waitGroup)
	go north.Northbound(&waitGroup, adminChannel)
	go reqBuilder.RequestBuilder(&waitGroup)
	go confMgr.ConfigManager(&waitGroup, adminChannel)
	go dataProcMgr.DataProcessingManager(&waitGroup)

	waitGroup.Wait()
}
