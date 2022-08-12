package deviceMonitor

import (
	"context"
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/proto/adapter"
	"github.com/onosproject/monitor-service/pkg/types"
)

var log = logger.GetLogger()

// TODO: Reduce duplicate code for update and shutdown commands, both have the same delete functionality

// Main function for managing monitoring of a single device
func DeviceMonitor(monitor types.DeviceMonitor) {
	var counterWaitGroup sync.WaitGroup
	var counterChannels []chan string

	// For each request create new goroutine dedicated to sending requests to adapter/device
	for index, req := range monitor.Requests {
		counterWaitGroup.Add(1)
		// Create a new channel and add it to list of channels, for communication directly to goroutine dedicated to sending requests
		counterChannels = append(counterChannels, make(chan string, 1))

		// Create new goroutine dedicated to sending requests to adapter/device
		go newCounter(req, monitor.DeviceName, monitor.Adapter, &counterWaitGroup, counterChannels[index])
	}

	// Create bool keeping track of state of the monitoring
	alive := true

	// While monitoring is active the loop continues
	for alive {
		// Block until command is sent on the managing channel
		cmd := <-monitor.ManagerChannel
		if cmd == "shutdown" {
			// Sends the shutdown command to all the goroutines dedicated to sending requests
			for index := 0; index < len(counterChannels); index++ {
				counterChannels[index] <- cmd
			}
			// Set monitoring to inactive
			alive = false
		} else if cmd == "update" {
			// Sends the shutdown command to all the goroutines dedicated to sending requests
			for index := 0; index < len(counterChannels); index++ {
				counterChannels[index] <- "shutdown"
			}

			// Block until new requests are sent on the request channel
			monitor.Requests = <-monitor.RequestsChannel

			// For each new request create a new goroutine dedicated to sending requests to adapter/device
			for index, req := range monitor.Requests {
				counterWaitGroup.Add(1)
				// Create a new channel and add it to list of channels, for communication directly to goroutine dedicated to sending requests
				counterChannels = append(counterChannels, make(chan string, 1))
				// Create new goroutine dedicated to sending requests to adapter/device
				go newCounter(req, monitor.DeviceName, monitor.Adapter, &counterWaitGroup, counterChannels[index])
			}
			log.Infof("Update complete for %v: %v\n", monitor.DeviceName, time.Now().UnixNano())
		}
	}

	// Wait until all goroutines dedicated to sending requests have stopped
	counterWaitGroup.Wait()
}

// Requests counters at the given interval, extract response and forward it
func newCounter(req types.Request, deviceName string, adapter *adapter.Adapter, waitGroup *sync.WaitGroup, counterChannel chan string) {
	defer waitGroup.Done()

	ctx := context.Background()

	// Create gNMI client
	c, err := createGnmiClient(adapter, ctx)
	if err != nil {
		// Restarts process after 10s, however, if the shutdown command is sent on
		// counterChannel, the process will stop
		restartTicker := time.NewTicker(10 * time.Second)

		select {
		case msg := <-counterChannel:
			if msg == "shutdown" {
				log.Info("Exits counter now")
				return
			}
		case <-restartTicker.C:
			// Retry creation of gNMI client (done by calling this function again)
			restartTicker.Stop()
			waitGroup.Add(1)
			go newCounter(req, deviceName, adapter, waitGroup, counterChannel)
			return
		}
	}

	counterIsActive := true
	id := 0

	// Create new goroutine that sends request to adapter/device
	go sendCounterReq(req, deviceName, ctx, c, &counterIsActive, id)

	// Start a ticker which will trigger repeatedly after (interval) milliseconds
	intervalTicker := time.NewTicker(time.Duration(req.Interval) * time.Millisecond)

	// Create new goroutine waiting for the repeated ticker
	go func() {
		for {
			select {
			case <-intervalTicker.C:
				// If no shutdown command has arrived yet, send "ticker" on "counterChannel"
				if counterIsActive {
					counterChannel <- "ticker"
				}
			default:
				if !counterIsActive {
					return
				}
			}
		}
	}()

	for counterIsActive {
		// Block until command/message arrive on "counterChannel"
		msg := <-counterChannel
		if msg == "shutdown" {
			// Stop ticker
			intervalTicker.Stop()
			// Set state of counter to inactive
			counterIsActive = false
		} else if msg == "ticker" {
			id += 1
			// Create new goroutine that sends request to adapter/device
			go sendCounterReq(req, deviceName, ctx, c, &counterIsActive, id)
		}
	}

	log.Infof("Exits %v from %v", req.Counters[0].Name, deviceName)
}
