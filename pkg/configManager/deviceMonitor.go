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

	// TEMPORARY (maybe not, depends if we need processing)
	"github.com/onosproject/monitor-service/pkg/streamManager"
)

// TODO: Place file in new folder representing its own module???

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

// Requests counters at the given interval, extract response and forward it.
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
			// Get the counter and send it to the data processing and .
			response, err := c.(*gclient.Client).Get(ctx, req.GnmiRequest)
			if err != nil {
				fmt.Printf("Target returned RPC error: %v", err)
			} else {
				// TODO: Send counter to data processing.

				// TODO: Use switch as name?
				extractData(response, req.GnmiRequest, "myOwnIdentifier" /*req.Name*/)
			}
		}
	}

	fmt.Println("Exits counter now")
}

// TODO: Parse all data, not just first notification to enable batch operations.
func extractData(response *gnmi.GetResponse, req *gnmi.GetRequest, name string) {
	// TODO: Rename adapterResponse to something like switchResponse.
	var adapterResponse types.AdapterResponse
	var schemaTree *types.SchemaTree

	if len(response.Notification) > 0 {

		if len(response.Notification[0].Update) == 0 {
			fmt.Printf("There is no data for request: %v\n", req)
			return
		}

		if err := proto.Unmarshal(response.Notification[0].Update[0].Val.GetProtoBytes(), &adapterResponse); err != nil {
			fmt.Printf("Failed to unmarshal ProtoBytes: %v", err)
		}

		// printTree(adapterResponse.Entries)

		// Get tree structure from slice.
		schemaTree = getTreeStructure(adapterResponse.Entries)

		// printTree(schemaTree)

		// Send data to subscription manager.

		// fmt.Println("---------schemaTree.Name---------")
		// fmt.Println(schemaTree.Name)
		// fmt.Println("------------req.Path-------------")
		// fmt.Println(req.Path)
		// fmt.Println("---------------------------------")

		sendDataToSubMgr(schemaTree, req.Path, name, adapterResponse.Timestamp)
		// addSchemaTreeValueToStream(schemaTree.Children[0], req.Path[0].Elem, 0, name, adapterResponse.Timestamp)
	}
}

func sendDataToSubMgr(schemaTree *types.SchemaTree, paths []*gnmi.Path, name string, adapterTs int64) {
	// Append values from the counters in the same order as the paths.
	var counterValues []string
	for index, counter := range schemaTree.Children {
		counterValues = append(counterValues, findCounterVal(counter, paths[index].Elem, 0, name, adapterTs))
	}

	fmt.Println("---------counterValues---------")
	fmt.Println(counterValues)
	fmt.Println("-------------Paths-------------")
	fmt.Println(paths)
	fmt.Println("-------------------------------")
}

func findCounterVal(schemaTree *types.SchemaTree, pathElems []*gnmi.PathElem, startIndex int, name string, adapterTs int64) string {
	fmt.Println("--------------")
	fmt.Printf("startIndex = %v\nlen(pathElems) = %v\n", startIndex, len(pathElems))
	if startIndex < len(pathElems) {
		fmt.Printf("pathElems[%v].Name = %v\nschemaTree.Name = %v\n", startIndex, pathElems[startIndex].Name, schemaTree.Name)
		if pathElems[startIndex].Name == schemaTree.Name {
			fmt.Printf("len(pathElems)-1 = %v\n", len(pathElems)-1)
			if startIndex == len(pathElems)-1 {
				return schemaTree.Value
			}
			for _, child := range schemaTree.Children {
				findCounterVal(child, pathElems, startIndex+1, name, adapterTs)
			}
		}
	}

	return ""
}

// Recursively find requested value(s) and send it to the subscription manager.
func addSchemaTreeValueToStream(schemaTree *types.SchemaTree, pathElems []*gnmi.PathElem, startIndex int, name string, adapterTs int64) {
	if startIndex < len(pathElems) {
		if pathElems[startIndex].Name == schemaTree.Name {
			if startIndex == len(pathElems)-1 {
				// fmt.Printf("Value sent: %v - %v\n", schemaTree.Name, schemaTree.Value)
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

func printTree(schemaTree *types.SchemaTree) {
	if schemaTree.Name != "" {
		fmt.Println(schemaTree.Name)
		if schemaTree.Value != "" {
			fmt.Println(schemaTree.Value)
		}

		for _, child := range schemaTree.Children {
			printTree(child)
		}
	}
}

// func printTree(schemaEntries []types.SchemaEntry) {
// 	for _, entry := range schemaEntries {
// 		fmt.Println(entry.Name)
// 		if entry.Value != "" {
// 			fmt.Println(entry.Value)
// 		}
// 	}
// }
