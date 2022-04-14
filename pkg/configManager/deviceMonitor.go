package deviceManager

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/onosproject/monitor-service/pkg/types"

	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
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
				extractData(response, r)
			}
		}
	}

	fmt.Println("Exits counter now")
}

func extractData(response *gnmi.GetResponse, req *gnmi.GetRequest) {
	var schema types.Schema
	var schemaTree *types.SchemaTree
	if len(response.Notification) > 0 {
		// Should replace serialization from json to proto, it is supposed to be faster.
		json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &schema)
		// This is not necessary if better serialization that can serialize recursive objects is used.
		schemaTree = getTreeStructure(schema)
	}

	// This is not necessary either if better serialization is used.
	// var val int
	// val, err = getSchemaTreeValue(schemaTree.Children[0], r.Path[0].Elem, 0)
	fmt.Printf("%s: ", req.Path[0].Target)
	getSchemaTreeValue(schemaTree.Children[0], req.Path[0].Elem, 0)

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(val)
	// }
}

func getSchemaTreeValue(schemaTree *types.SchemaTree, pathElems []*gnmi.PathElem, startIndex int) {
	if startIndex < len(pathElems) {
		if pathElems[startIndex].Name == schemaTree.Name {
			if startIndex == len(pathElems)-1 {
				// return strconv.Atoi(schemaTree.Value)
				fmt.Println(schemaTree.Value)
			}
			for _, child := range schemaTree.Children {
				getSchemaTreeValue(child, pathElems, startIndex+1)
			}
		}
	}

	// return -1, errors.New("Could not find value")
}

func getTreeStructure(schema types.Schema) *types.SchemaTree {
	var newTree *types.SchemaTree
	tree := &types.SchemaTree{}
	lastNode := ""
	for _, entry := range schema.Entries {
		if entry.Value == "" {
			// In a directory
			if entry.Tag == "end" {
				if entry.Name != "data" {
					if lastNode != "leaf" {
						// fmt.Println(tree.Name)
						tree = tree.Parent
					}
					lastNode = ""
					// continue
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
