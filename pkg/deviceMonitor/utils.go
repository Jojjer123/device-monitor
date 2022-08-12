package deviceMonitor

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"

	dataProcessing "github.com/onosproject/monitor-service/pkg/dataProcessingManager"
	"github.com/onosproject/monitor-service/pkg/proto/adapter"
	"github.com/onosproject/monitor-service/pkg/subscriptionManager"
	"github.com/onosproject/monitor-service/pkg/types"
)

// Creates a gNMI client
func createGnmiClient(adapter *adapter.Adapter, ctx context.Context) (client.Impl, error) {
	// Port should be a variable set earlier or should be changed to "10161" for secure communication
	port := "11161"

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{adapter.Address + ":" + port},
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		log.Errorf("Could not create a gNMI client: %+v", err)

		return nil, err
	}

	return c, nil
}

// Sends a request to adapter/device
func sendCounterReq(req types.Request, deviceName string, ctx context.Context, c client.Impl, active *bool, id int) {
	// Get the counter from the adapter/device
	response, err := c.(*gclient.Client).Get(ctx, req.GnmiRequest)

	if *active {
		if err != nil {
			log.Errorf("Target returned RPC error: %v", err)
		} else {
			// Extract data from the response
			extractData(response, req.GnmiRequest, deviceName)
		}
	}
}

// Extracts data from a gNMI response
func extractData(response *gnmi.GetResponse, req *gnmi.GetRequest, name string) {
	// Start of response reading where the gNMI response contains all paths to the values (without any extra serialization such as "adapterResponse")
	/*
		for _, notification := range response.Notification {
			if len(notification.GetUpdate()) == 0 {
				log.Warnf("No data in notification: %v", notification)
				continue
			}

			// TODO: extract data from update(s)
			for _, update := range notification.Update {
				// TODO: map val to path in an easily accessible format
				val, err := getValue(update)
				if err != nil {
					log.Errorf("Failed getting value from update: %v", err)
					continue
				}

				log.Infof("Mapping val: %s to path: %v", val, update.Path)
			}

			// TODO: pass data to data processing and subscription manager

		}
	*/

	var adapterResponse types.AdapterResponse
	var schemaTree *types.SchemaTree

	if len(response.Notification) > 0 {

		if len(response.Notification[0].Update) == 0 {
			log.Warnf("There is no data for request: %v", req)
			return
		}

		// Deserialize response that was built in adapter
		if err := proto.Unmarshal(response.Notification[0].Update[0].Val.GetProtoBytes(), &adapterResponse); err != nil {
			log.Errorf("Failed to unmarshal ProtoBytes: %v", err)
		}

		// Get tree structure of response
		schemaTree = getTreeStructure(adapterResponse.Entries)

		// Process data
		dataProcessing.ProcessData(schemaTree, req.Path)
		// Send data to subscription manager
		sendDataToSubMgr(schemaTree, req.Path, name, adapterResponse.Timestamp)
	}
}

// Takes in a gnmi.Update and converts the value to a string
func getValue(update *gnmi.Update) (string, error) {
	var value string

	// Get any kind of value, not just decimal values.
	switch update.Val.Value.(type) {
	case *gnmi.TypedValue_AnyVal:
		value = update.GetVal().GetAnyVal().String()
	case *gnmi.TypedValue_AsciiVal:
		value = update.GetVal().GetAsciiVal()
	case *gnmi.TypedValue_BoolVal:
		if update.GetVal().GetBoolVal() {
			value = "true"
		} else {
			value = "false"
		}
	case *gnmi.TypedValue_BytesVal:
		value = string(update.GetVal().GetBytesVal())
	case *gnmi.TypedValue_FloatVal:
		value = fmt.Sprintf("%f", update.GetVal().GetFloatVal())
	case *gnmi.TypedValue_DecimalVal:
		value = strconv.FormatInt(update.GetVal().GetDecimalVal().GetDigits(), 10)
	case *gnmi.TypedValue_IntVal:
		value = strconv.Itoa(int(update.GetVal().GetIntVal()))
	case *gnmi.TypedValue_JsonIetfVal:
		value = string(update.GetVal().GetJsonIetfVal())
	case *gnmi.TypedValue_JsonVal:
		value = string(update.GetVal().GetJsonVal())
	case *gnmi.TypedValue_LeaflistVal:
		value = update.GetVal().GetLeaflistVal().String()
	case *gnmi.TypedValue_ProtoBytes:
		value = string(update.GetVal().GetProtoBytes())
	case *gnmi.TypedValue_StringVal:
		value = update.GetVal().GetStringVal()
	case *gnmi.TypedValue_UintVal:
		value = strconv.FormatUint(update.GetVal().GetUintVal(), 10)
	default:
		log.Errorf("Value \"%v\" is not defined", update.GetVal())
		return "", errors.New("value not defined")
	}

	return value, nil
}

// Sends data to subscription manager
func sendDataToSubMgr(schemaTree *types.SchemaTree, paths []*gnmi.Path, name string, adapterTs int64) {
	// Append values from the counters in the same order as the paths.
	var counterValues []string

	// Find all the counter values
	for _, counter := range schemaTree.Children {
		findCounterVals(counter, &counterValues)
	}

	if len(counterValues) != len(paths) {
		log.Errorf("Failed to map counter values to paths with counters: %v\npaths: %v", counterValues, paths)
		return
	}

	// Send dictionary of paths mapped to counter values, to the subscription manager
	subscriptionManager.AddDataToSubscribers(createDictionary(counterValues, paths), name, adapterTs)
}

// Create a dictionary from the provided counter values and paths
func createDictionary(counterValues []string, paths []*gnmi.Path) []types.Dictionary {
	var dict []types.Dictionary

	for index, counterVal := range counterValues {
		dict = append(dict, types.Dictionary{
			paths[index].Elem[len(paths[index].Elem)-1].Name: counterVal,
		})
	}

	return dict
}

// Recursively finds the counter values from the schemaTree
func findCounterVals(schemaTree *types.SchemaTree, counterValues *[]string) {
	if schemaTree.Value != "" {
		// Check if all children of parent has values, then I must be a counter, otherwise I am just an identifier
		isIdentifier := false
		for _, child := range schemaTree.Parent.Children {
			// If child is directory
			if child.Value == "" {
				isIdentifier = true
			}
		}

		if !isIdentifier {
			*counterValues = append(*counterValues, schemaTree.Value)
		}
	} else {
		// Current schemaTree is directory
		for _, child := range schemaTree.Children {
			findCounterVals(child, counterValues)
		}
	}
}

// Build a schemaTree from the entries provided
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
