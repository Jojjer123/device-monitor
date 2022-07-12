package deviceMonitor

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"

	dataProcessing "github.com/onosproject/monitor-service/pkg/dataProcessingManager"
	"github.com/onosproject/monitor-service/pkg/proto/adapter"
	"github.com/onosproject/monitor-service/pkg/subscriptionManager"
	"github.com/onosproject/monitor-service/pkg/types"
)

func createGnmiClient(adapter *adapter.Adapter, target string, ctx context.Context) (client.Impl, error) {
	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{adapter.Address},
		Target:      target,
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		log.Errorf("Could not create a gNMI client: %v", err)

		return nil, err
	}

	return c, nil
}

func sendCounterReq(req types.Request, deviceName string, ctx context.Context, c client.Impl, active *bool, id int) {
	// fmt.Printf("Len of counter channel is: %v\n", len(counterChannel))

	fmt.Printf("Get %v from %v, req ID: %v, : %v\n", req.Counters[0].Name, deviceName, id, time.Now().UnixNano())

	// Get the counter and send it to the data processing and to possible subscribers.
	response, err := c.(*gclient.Client).Get(ctx, req.GnmiRequest)

	if *active {
		fmt.Printf("Received %v from %v, req ID: %v, : %v\n", req.Counters[0].Name, deviceName, id, time.Now().UnixNano())

		if err != nil {
			log.Errorf("Target returned RPC error: %v", err)
		} else {
			extractData(response, req.GnmiRequest, deviceName)
		}
	}
}

func extractData(response *gnmi.GetResponse, req *gnmi.GetRequest, name string) {
	// TODO: Rename adapterResponse to something like switchResponse.
	var adapterResponse types.AdapterResponse
	var schemaTree *types.SchemaTree

	if len(response.Notification) > 0 {

		if len(response.Notification[0].Update) == 0 {
			log.Warnf("There is no data for request: %v", req)
			return
		}

		if err := proto.Unmarshal(response.Notification[0].Update[0].Val.GetProtoBytes(), &adapterResponse); err != nil {
			log.Errorf("Failed to unmarshal ProtoBytes: %v", err)
		}

		// logger.Infof("Response entries: %v", adapterResponse.Entries)

		// Get tree structure from slice.
		schemaTree = getTreeStructure(adapterResponse.Entries)

		dataProcessing.ProcessData(schemaTree, req.Path)
		sendDataToSubMgr(schemaTree, req.Path, name, adapterResponse.Timestamp)
	}
}

func sendDataToSubMgr(schemaTree *types.SchemaTree, paths []*gnmi.Path, name string, adapterTs int64) {
	// Append values from the counters in the same order as the paths.
	// var counterValues []string
	var counterValues []string
	// logger.Infof("Number of children for schemaTree = %v", len(schemaTree.Children))
	// logger.Infof("SchemaTree child is: %v", schemaTree.Children[0])

	// for index, counter := range schemaTree.Children {
	for _, counter := range schemaTree.Children {
		// counterValues = append(counterValues, findCounterVal(counter, paths[index].Elem, 0))
		findCounterVals(counter, &counterValues)
	}

	// logger.Infof("Counter values: %v", counterValues)

	if len(counterValues) != len(paths) {
		log.Errorf("Failed to map counter values to paths with counters: %v\npaths: %v", counterValues, paths)
		return
	}

	// logger.Infof("Identifier %v is now calling the AddDataToSubscribers", name)
	subscriptionManager.AddDataToSubscribers(createDictionary(counterValues, paths), name, adapterTs)
}

func createDictionary(counterValues []string, paths []*gnmi.Path) []types.Dictionary {
	var dict []types.Dictionary

	for index, counterVal := range counterValues {
		dict = append(dict, types.Dictionary{
			paths[index].Elem[len(paths[index].Elem)-1].Name: counterVal,
		})
	}

	return dict
}

func findCounterVals(schemaTree *types.SchemaTree, counterValues *[]string) {
	if schemaTree.Value != "" {
		// Check if all children of parent has values, then I must be a counter, otherwise I am just an identifier.
		isIdentifier := false
		for _, child := range schemaTree.Parent.Children {
			// If child is directory.
			if child.Value == "" {
				isIdentifier = true
			}
		}

		if !isIdentifier {
			*counterValues = append(*counterValues, schemaTree.Value)
		}
	} else {
		// Current schemaTree is directory.
		for _, child := range schemaTree.Children {
			findCounterVals(child, counterValues)
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
