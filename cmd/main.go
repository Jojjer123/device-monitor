package main

import (
	"sync"

	confMgr "github.com/onosproject/monitor-service/pkg/configManager"
	dataProcMgr "github.com/onosproject/monitor-service/pkg/dataProcessingManager"
	north "github.com/onosproject/monitor-service/pkg/northbound"
	reqBuilder "github.com/onosproject/monitor-service/pkg/requestBuilder"
	storage "github.com/onosproject/monitor-service/pkg/storage"
	streamMgr "github.com/onosproject/monitor-service/pkg/streamManager"

	"github.com/onosproject/monitor-service/pkg/types"
)

const numberOfModules = 6

// Starts the main components of the monitor-service
func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfModules)

	// WARNING!!! potential problem: buffered vs unbuffered channels block in different stages of the communication.
	configAdminChannel := make(chan types.ConfigAdminChannelMessage)
	streamMgrChannel := make(chan types.StreamMgrChannelMessage)

	go storage.ConfigInterface(&waitGroup)
	go north.Northbound(&waitGroup, configAdminChannel, streamMgrChannel)
	go reqBuilder.RequestBuilder(&waitGroup)
	go confMgr.ConfigManager(&waitGroup, configAdminChannel)
	go dataProcMgr.DataProcessingManager(&waitGroup)
	go streamMgr.StreamManager(&waitGroup, streamMgrChannel)

	waitGroup.Wait()
}
