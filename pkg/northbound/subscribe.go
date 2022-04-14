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

	subRequest, err := stream.Recv()

	if err != nil {
		fmt.Print("Failed to receive from stream: ")
		fmt.Println(err)

	}

	// fmt.Print("target: ")
	// fmt.Println(subRequest.GetSubscribe().Prefix.Target)
	target := subRequest.GetSubscribe().Prefix.Target

	// fmt.Print("subscription: ")
	// fmt.Println(subRequest.GetSubscribe().GetSubscription())

	// fmt.Println(subRequest)

	// upd := []*pb.Update{
	// 	Value: &pb.Value{
	// 		Value: []byte("Hello from the other side"),
	// 	},
	// }

	// var notification *pb.Notification

	// notification.Update[0].Value.Value = []byte("test")

	var updateList []*pb.Update

	data := pb.TypedValue_StringVal{
		StringVal: "123Testing123",
	}

	path := pb.Path{
		Target: target,
		Elem: []*pb.PathElem{
			{
				Name: "interfaces",
			},
		},
	}

	update := pb.Update{
		Path: &path,
		Val: &pb.TypedValue{
			Value: &data,
		},
	}

	updateList = append(updateList, &update)

	response := pb.SubscribeResponse{
		Response: &pb.SubscribeResponse_Update{
			Update: &pb.Notification{
				Update: updateList,
			},
		},
	}

	for i := 0; i < 10; i++ {
		stream.Send(&response)
		time.Sleep(5 * time.Second)
	}

	return nil //s.Server.Subscribe(stream)
}
