package configManager

import (
	"fmt"

	"github.com/onosproject/monitor-service/pkg/deviceMonitor"
	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/proto/adapter"

	reqBuilder "github.com/onosproject/monitor-service/pkg/requestBuilder"
	"github.com/onosproject/monitor-service/pkg/types"
)

var log = logger.GetLogger()
var deviceMonitorStore []types.DeviceMonitor

func ExecuteAdminSetCmd(cmd string, target string, configIndex ...int) string {
	// Check exists as configIndex is not a required parameter, however, when used it should only be one integer
	if len(configIndex) > 1 {
		log.Warn("Config index should not be an array larger than 1")
	}

	switch cmd {
	case "Start":
		// Get slice of the different paths with their intervals (requests) and the appropriate adapter if one is necessary
		requests, adapter, deviceName := reqBuilder.GetRequestConf(target, configIndex[0])
		if len(requests) == 0 {
			return "No configurations to monitor"
		}
		startMonitoring(requests, adapter, target, deviceName)
	case "Update":
		// Get new paths and their intervals (requests)
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

// Stops all monitoring of a given device (target)
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

// Stops all monitoring of a device, then start monitor the same device with the updated requests
func updateMonitoring(requests []types.Request, target string) {
	for _, monitor := range deviceMonitorStore {
		// If target is found, send update command followed by new requests to the monitor goroutine
		if monitor.Target == target {
			fmt.Printf("Sending update to %v\n", monitor.Target)
			monitor.ManagerChannel <- "update"
			monitor.RequestsChannel <- requests
			return
		}
	}

	log.Warn("Could not find device monitor in store")
}

// Start monitoring the requests for the given device (target)
func startMonitoring(requests []types.Request, adapter *adapter.Adapter, target string, deviceName string) {
	// Consider checking Requests to update only if changed.
	monitor := types.DeviceMonitor{
		DeviceName:      deviceName,
		Target:          target,
		Adapter:         adapter,
		Requests:        requests,
		RequestsChannel: make(chan []types.Request, 1),
		ManagerChannel:  make(chan string, 1),
	}

	// Add monitor data object to global list, to keep track of all devices being monitored
	deviceMonitorStore = append(deviceMonitorStore, monitor)

	// Create a new goroutine with all required data for monitoring a device
	go deviceMonitor.DeviceMonitor(monitor)
}
