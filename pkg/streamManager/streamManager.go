package streamManager

import (
	"fmt"
	"sync"
	"time"

	"encoding/json"

	"github.com/onosproject/monitor-service/pkg/types"

	// "github.com/openconfig/gnmi/ctree"
	"github.com/openconfig/gnmi/proto/gnmi"
	// "github.com/openconfig/goyang/pkg/yang"
)

var streamStore []types.Stream

// TODO: Rename module to Subscription Manager

func StreamManager(waitGroup *sync.WaitGroup, streamMgrChannel chan types.StreamMgrChannelMessage) { //, adminChannel chan types.AdminChannelMessage) {
	// fmt.Println("Started StreamManager")
	defer waitGroup.Done()

	// TODO: Remove streamWaitGroup and add better way of keeping module "alive".

	var streamWaitGroup sync.WaitGroup

	// fmt.Println("Going to send function from StreamManager")
	var streamMgrMessage types.StreamMgrChannelMessage
	streamMgrMessage.ManageCmd = streamMgrCmd
	streamMgrChannel <- streamMgrMessage
	// fmt.Println("Sent function from StreamManager")

	streamWaitGroup.Wait()
	// fmt.Println("Closed StreamManager")
}

func streamMgrCmd(stream types.Stream, cmd string) string {
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
		}

	default:
		fmt.Printf("Did not recognize cmd: %s\n", cmd)
	}

	return ""
}

func AddDataToStream(dataVal string, subscriptionIdentifier string, adapterTs int64) types.Stream {
	for _, stream := range streamStore {
		if stream.Target[0].Name == subscriptionIdentifier {
			objectToSend := types.GatewayData{
				Data:             dataVal,
				MonitorTimestamp: time.Now().UnixNano(),
				AdapterTimestamp: adapterTs,
			}

			jsonBytes, err := json.Marshal(objectToSend)
			if err != nil {
				fmt.Printf("Failed to marshal to json, err: %v", err)
			}

			stream.StreamHandle.Send(&gnmi.SubscribeResponse{
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
			})
		}
	}

	return types.Stream{}
}
