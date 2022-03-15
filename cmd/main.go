package main

import (
	"sync"

	conf "github.com/onosproject/device-monitor/pkg/config"
	dataProc "github.com/onosproject/device-monitor/pkg/dataProcessing"
	deviceMgr "github.com/onosproject/device-monitor/pkg/deviceManager"
	north "github.com/onosproject/device-monitor/pkg/northbound"
	reqBuilder "github.com/onosproject/device-monitor/pkg/requestBuilder"
	topo "github.com/onosproject/device-monitor/pkg/topo"

	types "github.com/onosproject/device-monitor/pkg/types"
)

const numberOfComponents = 6

// Starts the main components of the device-monitor
func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfComponents)

	// WARNING potential problem: buffered vs unbuffered channels block in different stages of the communication.
	adminChannel := make(chan types.AdminChannelMessage)

	go topo.TopoInterface(&waitGroup)
	go conf.ConfigInterface(&waitGroup)
	go north.Northbound(&waitGroup, adminChannel)
	go reqBuilder.RequestBuilder(&waitGroup)
	go deviceMgr.DeviceManager(&waitGroup, adminChannel)
	go dataProc.DataProcessing(&waitGroup)

	waitGroup.Wait()
}
