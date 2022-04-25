package streamManager

import (
	"fmt"
	"sync"

	"github.com/onosproject/monitor-service/pkg/types"
	// "github.com/openconfig/gnmi/ctree"
	"github.com/openconfig/gnmi/proto/gnmi"
	// "github.com/openconfig/goyang/pkg/yang"
	// "google.golang.org/protobuf/proto"
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

// Needs to be renamed to something like AddDataToStream
func GetSubscriberStream(target string) types.Stream {
	// TODO: Add search for stream given the target.
	// var test types.Stream
	for index, stream := range streamStore {
		if index == 0 {
			test := stream

			// entry := yang.Entry{
			// 	Name:    "FirstEntry",
			// 	Kind:    yang.LeafEntry,
			// 	Default: "FirstVal",
			// }

			// tree := ctree.Tree{}
			// tree.Add([]string{"interface"}, entry)

			// bytesTree, err := proto.Marshal(tree)
			// if err != nil {
			// 	fmt.Printf("Failed to marshal tree with err: %v\n", err)
			// }

			stream.StreamHandle.Send(&gnmi.SubscribeResponse{
				Response: &gnmi.SubscribeResponse_Update{
					Update: &gnmi.Notification{
						Update: []*gnmi.Update{
							{
								Path: &gnmi.Path{
									Elem: stream.Target,
								},
								// Val: &gnmi.TypedValue{
								// 	Value: &gnmi.TypedValue_ProtoBytes{
								// 		ProtoBytes: bytesTree,
								// 	},
								// },
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

			fmt.Println(&gnmi.Path{
				Elem: test.Target,
			})
		}
	}

	return types.Stream{}
}
