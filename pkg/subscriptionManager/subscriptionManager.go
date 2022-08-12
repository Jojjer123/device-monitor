package subscriptionManager

import (
	"encoding/json"
	"time"

	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/types"
	"github.com/openconfig/gnmi/proto/gnmi"
)

var log = logger.GetLogger()
var streamStore []types.Stream

// Add or remove the stream provided (manages subscriptions)
func SubscriptionMgrCmd(stream types.Stream, cmd string) string {
	switch cmd {
	case "Add":
		streamStore = append(streamStore, stream)
	case "Remove":
		indexToBeRemoved := -1
		// Find index to remove
		for index, activeStream := range streamStore {
			if activeStream.StreamHandle == stream.StreamHandle {
				indexToBeRemoved = index
			}
		}

		if indexToBeRemoved != -1 {
			streamStore = append(streamStore[:indexToBeRemoved], streamStore[indexToBeRemoved+1:]...)
		} else {
			log.Warn("Could not find stream to delete")
		}

	default:
		log.Errorf("Did not recognize cmd: %s", cmd)
	}

	return ""
}

// Not tested for multiple subscribers to the same data
// Find subscribers to the data, and send the data to its subscribers
func AddDataToSubscribers(dataVal []types.Dictionary, subscriptionIdentifier string, adapterTs int64) {
	for _, stream := range streamStore {
		// If data has subscribers
		if stream.Target[0].Name == subscriptionIdentifier {
			objectToSend := types.GatewayData{
				Data:             dataVal,
				MonitorTimestamp: time.Now().UnixNano(),
				AdapterTimestamp: adapterTs,
			}

			// Serialize data using JSON
			jsonBytes, err := json.Marshal(objectToSend)
			if err != nil {
				log.Errorf("Failed to marshal to json, err: %v", err)
			}

			// Build subscribe response
			subResponse := &gnmi.SubscribeResponse{
				Response: &gnmi.SubscribeResponse_Update{
					Update: &gnmi.Notification{
						Timestamp: time.Now().UnixNano(),
						Update: []*gnmi.Update{
							{
								Path: &gnmi.Path{
									Elem: stream.Target,
								},
								Val: &gnmi.TypedValue{
									Value: &gnmi.TypedValue_JsonVal{
										JsonVal: jsonBytes,
									},
								},
							},
						},
					},
				},
			}

			log.Infof("Send data from %v, to gnmi-gateway: %v\n", subscriptionIdentifier, time.Now().UnixNano())

			// Send subscribe response with the data
			stream.StreamHandle.Send(subResponse)
		}
	}
}
