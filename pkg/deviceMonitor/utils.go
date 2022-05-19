package deviceMonitor

import (
	"context"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"

	dataProcessing "github.com/onosproject/monitor-service/pkg/dataProcessingManager"
	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/subscriptionManager"
	"github.com/onosproject/monitor-service/pkg/types"
)

func createGnmiClient(adapter types.Adapter, target string, ctx context.Context) (client.Impl, error) {
	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{adapter.Address},
		Target:      target,
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		logger.Errorf("Could not create a gNMI client: %v", err)

		return nil, err
	}

	return c, nil
}

func extractData(response *gnmi.GetResponse, req *gnmi.GetRequest, name string) {
	// TODO: Rename adapterResponse to something like switchResponse.
	var adapterResponse types.AdapterResponse
	var schemaTree *types.SchemaTree

	if len(response.Notification) > 0 {

		if len(response.Notification[0].Update) == 0 {
			logger.Warnf("There is no data for request: %v", req)
			return
		}

		if err := proto.Unmarshal(response.Notification[0].Update[0].Val.GetProtoBytes(), &adapterResponse); err != nil {
			logger.Errorf("Failed to unmarshal ProtoBytes: %v", err)
		}

		// Get tree structure from slice.
		schemaTree = getTreeStructure(adapterResponse.Entries)

		dataProcessing.ProcessData(schemaTree, req.Path)
		sendDataToSubMgr(schemaTree, req.Path, name, adapterResponse.Timestamp)
	}
}

func sendDataToSubMgr(schemaTree *types.SchemaTree, paths []*gnmi.Path, name string, adapterTs int64) {
	// Append values from the counters in the same order as the paths.
	var counterValues []string
	for index, counter := range schemaTree.Children {
		counterValues = append(counterValues, findCounterVal(counter, paths[index].Elem, 0))
	}

	if len(counterValues) != len(paths) {
		logger.Error("Failed to map counter values to paths.")
		return
	}

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

// Call findCounterVal with startIndex as 0, in order to start searching through pathElems from index 0.
func findCounterVal(schemaTree *types.SchemaTree, pathElems []*gnmi.PathElem, startIndex int) string {
	if startIndex < len(pathElems) {
		if pathElems[startIndex].Name == schemaTree.Name {
			if startIndex == len(pathElems)-1 {
				return schemaTree.Value
			}
			var childResult string
			for _, child := range schemaTree.Children {
				childResult += findCounterVal(child, pathElems, startIndex+1)
			}
			return childResult
		}
	}

	return ""
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
