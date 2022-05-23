package deviceMonitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/types"
	gclient "github.com/openconfig/gnmi/client/gnmi"
)

// TODO: Place file in new folder representing its own module???
// TODO: Split this file into at least one more, for some helpers.

func DeviceMonitor(monitor types.DeviceMonitor) {
	var counterWaitGroup sync.WaitGroup
	var counterChannels []chan string

	// fmt.Printf("Requests: %v\n", monitor.Requests)

	for index, req := range monitor.Requests {
		counterWaitGroup.Add(1)
		counterChannels = append(counterChannels, make(chan string, 1))

		// fmt.Printf("Sending channel %v to %v\n", index, req.Counters[0].Name)

		go newCounter(req, monitor.DeviceName, monitor.Target, monitor.Adapter, &counterWaitGroup, counterChannels[index])
	}

	// fmt.Println(len(monitor.Requests))
	// fmt.Println(len(counterChannels))

	alive := true
	for alive {
		cmd := <-monitor.ManagerChannel
		if cmd == "shutdown" {
			// fmt.Printf("len: %v\n", len(counterChannels))
			// fmt.Printf("Shutting down %v:\n", monitor.Target)
			for index := 0; index < len(counterChannels); index++ {
				// fmt.Println(index)
				// fmt.Println(cap(counterChannels[index]))
				counterChannels[index] <- cmd
				// fmt.Println("Sent command on channel now")
			}
			alive = false
		} else if cmd == "update" {
			for index := 0; index < len(counterChannels); index++ {
				counterChannels[index] <- "shutdown"
			}

			monitor.Requests = <-monitor.RequestsChannel

			for index, req := range monitor.Requests {
				counterWaitGroup.Add(1)
				counterChannels = append(counterChannels, make(chan string, 1))
				go newCounter(req, monitor.DeviceName, monitor.Target, monitor.Adapter, &counterWaitGroup, counterChannels[index])
			}
		}
	}

	counterWaitGroup.Wait()
}

// Requests counters at the given interval, extract response and forward it.
func newCounter(req types.Request, deviceName string, target string, adapter types.Adapter, waitGroup *sync.WaitGroup, counterChannel chan string) {
	defer waitGroup.Done()

	ctx := context.Background()

	c, err := createGnmiClient(adapter, target, ctx)
	if err != nil {
		// Restarts process after 10s, however, if the shutdown command is sent on
		// counterChannel, the process will stop.

		restartTicker := time.NewTicker(10 * time.Second)

		select {
		case msg := <-counterChannel:
			if msg == "shutdown" {
				logger.Info("Exits counter now")
				return
			}
		// case <-time.After(10 * time.Second):
		case <-restartTicker.C:
			restartTicker.Stop()
			waitGroup.Add(1)
			go newCounter(req, deviceName, target, adapter, waitGroup, counterChannel)
			return
		}
	}

	fmt.Printf("Get %v from %v: %v\n", req.Counters[0].Name, deviceName, time.Now().UnixNano())

	// Get the counter and send it to the data processing and to possible subscribers.
	response, err := c.(*gclient.Client).Get(ctx, req.GnmiRequest)

	fmt.Printf("Received %v from %v: %v\n", req.Counters[0].Name, deviceName, time.Now().UnixNano())

	if err != nil {
		logger.Errorf("Target returned RPC error: %v", err)
	} else {
		extractData(response, req.GnmiRequest, deviceName)
	}

	// Start a ticker which will trigger repeatedly after (interval) milliseconds.
	intervalTicker := time.NewTicker(time.Duration(req.Interval) * time.Millisecond)

	counterIsActive := true

	go func() {
		select {
		case <-intervalTicker.C:
			if counterIsActive {
				counterChannel <- "ticker"
			}
		default:
			if !counterIsActive {
				return
			}
		}
	}()

	for counterIsActive {
		// select {
		// case msg := <-counterChannel:
		msg := <-counterChannel
		if msg == "shutdown" {
			// fmt.Println("Shutdown message arrived")
			intervalTicker.Stop()
			counterIsActive = false
		} else if msg == "ticker" {
			// fmt.Printf("Len of counter channel is: %v\n", len(counterChannel))

			fmt.Printf("Get %v from %v: %v\n", req.Counters[0].Name, deviceName, time.Now().UnixNano())

			// Get the counter and send it to the data processing and to possible subscribers.
			response, err := c.(*gclient.Client).Get(ctx, req.GnmiRequest)

			fmt.Printf("Received %v from %v: %v\n", req.Counters[0].Name, deviceName, time.Now().UnixNano())

			if err != nil {
				logger.Errorf("Target returned RPC error: %v", err)
			} else {
				extractData(response, req.GnmiRequest, deviceName)
			}
		}
		// } else {
		// 	logger.Errorf("Counter channel message is not \"shutdown\", it is: %v", msg)
		// }
		// case <-intervalTicker.C:
		// 	fmt.Printf("Len of counter channel is: %v\n", len(counterChannel))

		// 	fmt.Printf("Get %v from %v: %v\n", req.Counters[0].Name, deviceName, time.Now().UnixNano())

		// 	// Get the counter and send it to the data processing and to possible subscribers.
		// 	response, err := c.(*gclient.Client).Get(ctx, req.GnmiRequest)

		// 	fmt.Printf("Received %v from %v: %v\n", req.Counters[0].Name, deviceName, time.Now().UnixNano())

		// 	if err != nil {
		// 		logger.Errorf("Target returned RPC error: %v", err)
		// 	} else {
		// 		extractData(response, req.GnmiRequest, deviceName)
		// 	}
		// default:
		// }
	}

	logger.Infof("Exits %v from %v", req.Counters[0].Name, deviceName)
}
