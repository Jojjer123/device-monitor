package main

import (
	"sync"

	adm "github.com/onosproject/device-monitor/pkg/admin"
	conf "github.com/onosproject/device-monitor/pkg/config"
	dataProc "github.com/onosproject/device-monitor/pkg/dataProcessing"
	deviceMgr "github.com/onosproject/device-monitor/pkg/deviceManager"
	kafka "github.com/onosproject/device-monitor/pkg/kafka"
	reqBuilder "github.com/onosproject/device-monitor/pkg/requestBuilder"
	topo "github.com/onosproject/device-monitor/pkg/topo"
)

const numberOfComponents = 7

// Starts the main components of the device-monitor
func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfComponents)

	// WARNING potential problem: buffered vs unbuffered channels block in different stages of the communication.
	adminChannel := make(chan string)

	go topo.TopoInterface(&waitGroup)
	go conf.ConfigInterface(&waitGroup)
	go adm.AdminInterface(&waitGroup, adminChannel)
	go kafka.KafkaInterface(&waitGroup)
	go reqBuilder.RequestBuilder(&waitGroup)
	go deviceMgr.DeviceManager(&waitGroup, adminChannel)
	go dataProc.DataProcessing(&waitGroup)

	waitGroup.Wait()
}
