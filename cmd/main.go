package main

import (
	"sync"

	confMgr "github.com/onosproject/monitor-service/pkg/configManager"
	dataProcMgr "github.com/onosproject/monitor-service/pkg/dataProcessingManager"
	north "github.com/onosproject/monitor-service/pkg/northbound"
	reqBuilder "github.com/onosproject/monitor-service/pkg/requestBuilder"
	streamMgr "github.com/onosproject/monitor-service/pkg/streamManager"

	"github.com/onosproject/monitor-service/pkg/types"
)

const numberOfModules = 5

// Starts the main components of the monitor-service
func main() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(numberOfModules)

	configAdminChannel := make(chan types.ConfigAdminChannelMessage)
	streamMgrChannel := make(chan types.StreamMgrChannelMessage)

	go north.Northbound(&waitGroup, configAdminChannel, streamMgrChannel)
	go reqBuilder.RequestBuilder(&waitGroup)
	go confMgr.ConfigManager(&waitGroup, configAdminChannel)
	go dataProcMgr.DataProcessingManager(&waitGroup)
	go streamMgr.StreamManager(&waitGroup, streamMgrChannel)

	waitGroup.Wait()
}
