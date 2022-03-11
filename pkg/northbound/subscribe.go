package northboundInterface

import (
	"fmt"
	"time"

	"github.com/google/gnxi/utils/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/openconfig/gnmi/proto/gnmi"
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

	message, err := stream.Recv()

	if err != nil {
		fmt.Print("Failed to receive from stream: ")
		fmt.Println(err)

	}

	fmt.Println(message)

	// upd := []*pb.Update{
	// 	Value: &pb.Value{
	// 		Value: []byte("Hello from the other side"),
	// 	},
	// }

	// var notification *pb.Notification

	// notification.Update[0].Value.Value = []byte("test")

	var upd []*pb.Update

	fuckingHell := pb.TypedValue_StringVal{
		StringVal: "from the other fucking side",
	}

	reee := pb.Update{
		Val: &pb.TypedValue{
			Value: &fuckingHell,
		},
	}

	upd = append(upd, &reee)

	response := pb.SubscribeResponse{
		Response: &pb.SubscribeResponse_Update{
			Update: &pb.Notification{
				Update: upd,
			},
		},
	}
	stream.Send(&response)
	time.Sleep(5 * time.Second)
	stream.Send(&response)

	return nil //s.Server.Subscribe(stream)
}
