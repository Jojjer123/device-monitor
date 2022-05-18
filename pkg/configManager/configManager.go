package configManager

import (
	"github.com/onosproject/monitor-service/pkg/deviceMonitor"
	"github.com/onosproject/monitor-service/pkg/logger"
	reqBuilder "github.com/onosproject/monitor-service/pkg/requestBuilder"
	"github.com/onosproject/monitor-service/pkg/types"
)

var deviceMonitorStore []types.DeviceMonitor

func ExecuteAdminSetCmd(cmd string, target string, configIndex ...int) string {
	if len(configIndex) > 1 {
		logger.Warn("Config index should not be an array larger than 1")
	}

	switch cmd {
	case "Create":
		// Get slice of the different paths with their intervals and the appropriate adapter if one is necessary
		requests, adapter := reqBuilder.GetConfig(target, configIndex[0])
		createDeviceMonitor(requests, adapter, target)
	case "Update":
		requests, _ := reqBuilder.GetConfig(target, configIndex[0])
		updateDeviceMonitor(requests, target)
	case "Delete":
		deleteDeviceMonitor(target)
	default:
		logger.Warnf("Could not find command: %v", cmd)
		return "Could not find command: " + cmd
	}

	return "Successfully executed command sent"
}

func deleteDeviceMonitor(target string) {
	for index, monitor := range deviceMonitorStore {
		if monitor.Target == target {
			monitor.ManagerChannel <- "shutdown"

			deviceMonitorStore[index] = deviceMonitorStore[len(deviceMonitorStore)-1]
			deviceMonitorStore = deviceMonitorStore[:len(deviceMonitorStore)-1]
			return
		}
	}

	logger.Warn("Could not find device monitor in store")
}

func updateDeviceMonitor(requests []types.Request, target string) {
	for _, monitor := range deviceMonitorStore {
		if monitor.Target == target {
			monitor.ManagerChannel <- "update"
			monitor.RequestsChannel <- requests
			return
		}
	}

	logger.Warn("Could not find device monitor in store")
}

func createDeviceMonitor(requests []types.Request, adapter types.Adapter, target string) {
	// Consider checking Requests to update only if changed.
	monitor := types.DeviceMonitor{
		Target:          target,
		Adapter:         adapter,
		Requests:        requests,
		RequestsChannel: make(chan []types.Request),
		ManagerChannel:  make(chan string),
	}

	deviceMonitorStore = append(deviceMonitorStore, monitor)

	go deviceMonitor.DeviceMonitor(monitor)
}
