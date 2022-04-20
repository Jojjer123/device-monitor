package northboundInterface

import (
	"fmt"

	"github.com/google/gnxi/utils/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/openconfig/gnmi/proto/gnmi"

	"github.com/onosproject/monitor-service/pkg/types"
)

func (s *server) Subscribe(stream pb.GNMI_SubscribeServer) error {
	msg, ok := credentials.AuthorizeUser(stream.Context())
	if !ok {
		fmt.Print("Denied a Subscribe request: ")
		fmt.Println(msg)

		return status.Error(codes.PermissionDenied, msg)
	}

	fmt.Print("Allowed a Subscribe request: ")
	fmt.Println(msg)

	subRequest, err := stream.Recv()

	if err != nil {
		fmt.Print("Failed to receive from stream: ")
		fmt.Println(err)
	}

	// TODO: Add stream-handle to table of active subscriptions with the topic that the stream is subscribing to.

	newStream := types.Stream{
		StreamHandle: stream,
		Target:       subRequest.GetSubscribe().Prefix.Target,
	}

	s.StreamMgrCmd(newStream, "Add")

	// TODO: Add response notifying sender the success of creation?

	// target := subRequest.GetSubscribe().Prefix.Target

	// path := pb.Path{
	// 	Target: target,
	// 	Elem: []*pb.PathElem{
	// 		{
	// 			Name: "interfaces",
	// 		},
	// 	},
	// }

	// val := &pb.TypedValue{
	// 	Value: &pb.TypedValue_StringVal{
	// 		StringVal: "123Testing123",
	// 	},
	// }

	// response := pb.SubscribeResponse{
	// 	Response: &pb.SubscribeResponse_Update{
	// 		Update: &pb.Notification{
	// 			Update: []*pb.Update{
	// 				{
	// 					Path: &path,
	// 					Val:  val,
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// for i := 0; i < 10; i++ {
	// 	stream.Send(&response)
	// 	time.Sleep(5 * time.Second)
	// }

	return nil //s.Server.Subscribe(stream)
}
