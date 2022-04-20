package streamManager

import (
	"fmt"
	"sync"

	"github.com/onosproject/monitor-service/pkg/types"
	"github.com/openconfig/gnmi/proto/gnmi"
)

var streamStore []types.Stream

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
	// fmt.Println("Cmd arrived to StreamManager")
	switch cmd {
	case "Add":
		streamStore = append(streamStore, stream)
	default:
		fmt.Printf("Did not recognize cmd: %s\n", cmd)
	}
	fmt.Println(stream.Target, cmd)

	return ""
}

func GetSubscriberStream(target string) types.Stream {
	// TODO: Add search for stream given the target.
	var test types.Stream
	for index, stream := range streamStore {
		if index == 0 {
			test = stream
			stream.StreamHandle.Send(&gnmi.SubscribeResponse{
				Response: &gnmi.SubscribeResponse_Update{
					Update: &gnmi.Notification{
						Update: []*gnmi.Update{
							{
								Path: &gnmi.Path{
									Elem: stream.Target,
								},
								Val: &gnmi.TypedValue{
									Value: &gnmi.TypedValue_StringVal{
										StringVal: target,
									},
								},
							},
						},
					},
				},
			})
		}
	}

	fmt.Println(&gnmi.Path{
		Elem: test.Target,
	})

	return types.Stream{}
}
