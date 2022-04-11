package deviceManager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
				var schema Schema
				var schemaTree *SchemaTree
				if len(response.Notification) > 0 {
					// Should replace serialization from json to proto, it is supposed to be faster.
					json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &schema)
					// This is not necessary if better serialization that can serialize recursive objects is used.
					schemaTree = getTreeStructure(schema)
				}

				// This is not necessary either if better serialization is used.
				var val int
				val, err = getSchemaTreeValue(schemaTree.Children[0], r.Path[0].Elem, 0)

				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(val)
				}
			}
		}
	}

	fmt.Println("Exits counter now")
}

func getSchemaTreeValue(schemaTree *SchemaTree, pathElems []*gnmi.PathElem, startIndex int) (int, error) {
	if startIndex < len(pathElems) {
		if pathElems[startIndex].Name == schemaTree.Name {
			if startIndex == len(pathElems)-1 {
				return strconv.Atoi(schemaTree.Value)
			}
			for _, child := range schemaTree.Children {
				getSchemaTreeValue(child, pathElems, startIndex+1)
			}
		}
	}

	return -1, errors.New("Could not find value")
}

func printSchemaTree(schemaTree *SchemaTree) {
	fmt.Printf("%s - %s - %v - %s\n", schemaTree.Parent.Name, schemaTree.Name, schemaTree.Namespace, schemaTree.Value)
	for _, child := range schemaTree.Children {
		printSchemaTree(child)
	}
}

func getTreeStructure(schema Schema) *SchemaTree {
	var newTree *SchemaTree
	tree := &SchemaTree{}
	lastNode := ""
	for _, entry := range schema.Entries {
		if entry.Value == "" { // Directory
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

				newTree = &SchemaTree{Parent: tree}

				newTree.Name = entry.Name
				newTree.Namespace = entry.Namespace
				newTree.Parent.Children = append(newTree.Parent.Children, newTree)

				tree = newTree
			}
		} else { // Leaf
			newTree = &SchemaTree{Parent: tree}

			newTree.Name = entry.Name
			newTree.Value = entry.Value
			newTree.Parent.Children = append(newTree.Parent.Children, newTree)

			lastNode = "leaf"
		}
	}
	return tree
}

type SchemaTree struct {
	Name      string
	Namespace string
	Children  []*SchemaTree
	Parent    *SchemaTree
	Value     string
}

type Schema struct {
	Entries []SchemaEntry
}

type SchemaEntry struct {
	Name      string
	Tag       string
	Namespace string
	Value     string
}
