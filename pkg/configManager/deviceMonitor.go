package deviceManager

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/monitor-service/pkg/types"

	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"

	// TEMPORARY
	"github.com/onosproject/monitor-service/pkg/streamManager"
)

// target string, adapter types.Adapter, requests []types.Request, managerChannel <-chan string
func deviceMonitor(monitor types.DeviceMonitor) {
	var counterWaitGroup sync.WaitGroup
	var counterChannels []chan string

	for index, req := range monitor.Requests {
		counterWaitGroup.Add(1)
		counterChannels = append(counterChannels, make(chan string))
		go newCounter(req, monitor.Target, monitor.Adapter, &counterWaitGroup, counterChannels[index])
	}

	alive := true
	for alive {
		cmd := <-monitor.ManagerChannel
		if cmd == "shutdown" {
			for _, ch := range counterChannels {
				ch <- cmd
			}
			alive = false
		} else if cmd == "update" {
			for _, ch := range counterChannels {
				ch <- "shutdown"
			}

			monitor.Requests = <-monitor.RequestsChannel

			for index, req := range monitor.Requests {
				counterWaitGroup.Add(1)
				counterChannels = append(counterChannels, make(chan string))
				go newCounter(req, monitor.Target, monitor.Adapter, &counterWaitGroup, counterChannels[index])
			}
		}
	}

	counterWaitGroup.Wait()
}

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

	// Start a ticker which will trigger repeatedly after (interval) milliseconds.
	intervalTicker := time.NewTicker(time.Duration(req.Interval) * time.Millisecond)

	counterIsActive := true
	for counterIsActive {
		select {
		case msg := <-counterChannel:
			if msg == "shutdown" {
				intervalTicker.Stop()
				counterIsActive = false
			}
		case <-intervalTicker.C:
			// Get the counter here and send it to the data processing.
			response, err := c.(*gclient.Client).Get(ctx, r)
			if err != nil {
				fmt.Printf("Target returned RPC error: %v", err)
			} else {
				// TODO: Send counter to data processing.
				extractData(response, r, req.Name)
			}
		}
	}

	fmt.Println("Exits counter now")
}

func extractData(response *gnmi.GetResponse, req *gnmi.GetRequest, name string) {
	var adapterResponse types.AdapterResponse
	var schemaTree *types.SchemaTree
	if len(response.Notification) > 0 {
		// Should replace serialization from json to proto, it is supposed to be faster.
		if err := proto.Unmarshal(response.Notification[0].Update[0].Val.GetProtoBytes(), &adapterResponse); err != nil {
			fmt.Printf("Failed to unmarshal ProtoBytes: %v", err)
		}
		// json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &adapterResponse)

		startTime := time.Now().UnixNano()
		schemaTree = getTreeStructure(adapterResponse.Entries)
		fmt.Printf("Time to rebuild tree: %v", time.Now().UnixNano()-startTime)
	}

	addSchemaTreeValueToStream(schemaTree.Children[0], req.Path[0].Elem, 0, name, adapterResponse.Timestamp)
}

func addSchemaTreeValueToStream(schemaTree *types.SchemaTree, pathElems []*gnmi.PathElem, startIndex int, name string, adapterTs int64) {
	if startIndex < len(pathElems) {
		if pathElems[startIndex].Name == schemaTree.Name {
			if startIndex == len(pathElems)-1 {
				// fmt.Println(schemaTree.Value)
				streamManager.AddDataToStream(schemaTree.Value, name, adapterTs)
			}
			for _, child := range schemaTree.Children {
				addSchemaTreeValueToStream(child, pathElems, startIndex+1, name, adapterTs)
			}
		}
	}
}

func getTreeStructure(schemaEntries []types.SchemaEntry) *types.SchemaTree {
	var newTree *types.SchemaTree
	tree := &types.SchemaTree{}
	lastNode := ""
	for _, entry := range schemaEntries {
		if entry.Value == "" {
			// In a directory
			if entry.Tag == "end" {
				if entry.Name != "data" {
					if lastNode != "leaf" {
						tree = tree.Parent
					}
					lastNode = ""
				}
			} else {
				newTree = &types.SchemaTree{Parent: tree}

				newTree.Name = entry.Name
				newTree.Namespace = entry.Namespace
				newTree.Parent.Children = append(newTree.Parent.Children, newTree)

				tree = newTree
			}
		} else {
			// In a leaf
			newTree = &types.SchemaTree{Parent: tree}

			newTree.Name = entry.Name
			newTree.Value = entry.Value
			newTree.Parent.Children = append(newTree.Parent.Children, newTree)

			lastNode = "leaf"
		}
	}
	return tree
}
