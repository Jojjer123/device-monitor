package configManager

import (
	"fmt"

	"github.com/onosproject/monitor-service/pkg/deviceMonitor"
	"github.com/onosproject/monitor-service/pkg/logger"
	reqBuilder "github.com/onosproject/monitor-service/pkg/requestBuilder"
	"github.com/onosproject/monitor-service/pkg/types"
)

var log = logger.GetLogger()
var deviceMonitorStore []types.DeviceMonitor

func ExecuteAdminSetCmd(cmd string, target string, configIndex ...int) string {
	if len(configIndex) > 1 {
		log.Warn("Config index should not be an array larger than 1")
	}

	switch cmd {
	case "Start":
		// Get slice of the different paths with their intervals and the appropriate adapter if one is necessary
		// Should create new object with all the data inside.
		requests, adapter, deviceName := reqBuilder.GetRequestConf(target, configIndex[0])
		if len(requests) == 0 {
			return "No configurations to monitor"
		}
		startMonitoring(requests, adapter, target, deviceName)
	case "Update":
		requests, _, _ := reqBuilder.GetRequestConf(target, configIndex[0])
		if len(requests) == 0 {
			return "No configurations to monitor"
		}
		updateMonitoring(requests, target)
	case "Stop":
		stopMonitoring(target)
	default:
		log.Warnf("Could not find command: %v", cmd)
		return "Could not find command: " + cmd
	}

	return "Successfully executed command"
}

func stopMonitoring(target string) {
	for index, monitor := range deviceMonitorStore {
		if monitor.Target == target {
			monitor.ManagerChannel <- "shutdown"

			deviceMonitorStore[index] = deviceMonitorStore[len(deviceMonitorStore)-1]
			deviceMonitorStore = deviceMonitorStore[:len(deviceMonitorStore)-1]
			return
		}
	}

	log.Warn("Could not find device monitor in store")
}

func updateMonitoring(requests []types.Request, target string) {
	for _, monitor := range deviceMonitorStore {
		if monitor.Target == target {
			fmt.Printf("Sending update to %v\n", monitor.Target)
			monitor.ManagerChannel <- "update"
			monitor.RequestsChannel <- requests
			return
		}
	}

	log.Warn("Could not find device monitor in store")
}

func startMonitoring(requests []types.Request, adapter types.Adapter, target string, deviceName string) {
	// Consider checking Requests to update only if changed.
	monitor := types.DeviceMonitor{
		DeviceName:      deviceName,
		Target:          target,
		Adapter:         adapter,
		Requests:        requests,
		RequestsChannel: make(chan []types.Request, 1),
		ManagerChannel:  make(chan string, 1),
	}

	deviceMonitorStore = append(deviceMonitorStore, monitor)

	go deviceMonitor.DeviceMonitor(monitor)
}
