package deviceManager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/onosproject/device-monitor/pkg/types"

	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
)

func deviceMonitor(target string, adapter types.Adapter, requests []types.Request, managerChannel <-chan string) {
	var counterWaitGroup sync.WaitGroup
	var counterChannels []chan string

	for index, req := range requests {
		counterWaitGroup.Add(1)
		counterChannels = append(counterChannels, make(chan string))
		go newCounter(req, target, adapter, &counterWaitGroup, counterChannels[index])
	}

	alive := true
	for alive {
		x := <-managerChannel
		if x == "shutdown" {
			fmt.Println("Received shutdown command on channel now...")

			for _, ch := range counterChannels {
				ch <- x
			}
			alive = false
		}
	}

	fmt.Println("Shutting down device monitor now...")
	counterWaitGroup.Wait()
}

// req, target, adapter, counterWaitGroup, counterChannels[index]
func newCounter(req types.Request, target string, adapter types.Adapter, waitGroup *sync.WaitGroup, counterChannel <-chan string) {
	defer waitGroup.Done()

	ctx := context.Background()

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{adapter.Address},
		Target:      target,
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	r := &gnmi.GetRequest{
		Type: gnmi.GetRequest_STATE,
		Path: []*gnmi.Path{
			{
				Target: target,
				Elem:   req.Path,
			},
		},
	}

	// Start a ticker which will trigger repeatedly after (interval) seconds.
	intervalTicker := time.NewTicker(time.Duration(req.Interval) * time.Millisecond)

	counterIsActive := true
	for counterIsActive {
		select {
		case <-counterChannel:
			intervalTicker.Stop()
			counterIsActive = false
		case <-intervalTicker.C:
			// TODO: Get the counters for the given interval here and send them to the data processing part.
			// fmt.Println("Ticker triggered")
			response, err := c.(*gclient.Client).Get(ctx, r)
			if err != nil {
				fmt.Print("Target returned RPC error for Get(")
				fmt.Print(r.String())
				fmt.Print("): ")
				fmt.Println(err)
			} else {
				fmt.Println(response)
			}
		}
	}

	fmt.Println("Exits counter now")
}
