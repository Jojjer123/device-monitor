package subscriptionManager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/onosproject/monitor-service/pkg/logger"
	"github.com/onosproject/monitor-service/pkg/types"
	"github.com/openconfig/gnmi/proto/gnmi"
)

var streamStore []types.Stream

func SubscriptionMgrCmd(stream types.Stream, cmd string) string {
	switch cmd {
	case "Add":
		streamStore = append(streamStore, stream)
	case "Remove":
		indexToBeRemoved := -1
		for index, activeStream := range streamStore {
			if activeStream.StreamHandle == stream.StreamHandle {
				indexToBeRemoved = index
			}
		}

		if indexToBeRemoved != -1 {
			streamStore = append(streamStore[:indexToBeRemoved], streamStore[indexToBeRemoved+1:]...)
		} else {
			logger.Warn("Could not find stream to delete")
		}

	default:
		logger.Errorf("Did not recognize cmd: %s", cmd)
	}

	return ""
}

func AddDataToSubscribers(dataVal []types.Dictionary, subscriptionIdentifier string, adapterTs int64) {
	for _, stream := range streamStore {
		if stream.Target[0].Name == subscriptionIdentifier {
			objectToSend := types.GatewayData{
				Data:             dataVal,
				MonitorTimestamp: time.Now().UnixNano(),
				AdapterTimestamp: adapterTs,
			}

			jsonBytes, err := json.Marshal(objectToSend)
			if err != nil {
				logger.Errorf("Failed to marshal to json, err: %v", err)
			}

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

			fmt.Printf("Send data from %v, to gnmi-gateway: %v\n", subscriptionIdentifier, time.Now().UnixNano())

			stream.StreamHandle.Send(subResponse)
		}
	}
}
