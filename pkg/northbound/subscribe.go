package northboundInterface

import (
	"fmt"
	"time"

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

	// fmt.Println(subRequest.GetSubscribe().Subscription[0].Path)

	newStream := types.Stream{
		StreamHandle: stream,
		Target:       subRequest.GetSubscribe().Subscription[0].Path.Elem, // Previously was *&subRequest.GetSubscribe()...
	}

	s.StreamMgrCmd(newStream, "Add")

	go func() {
		_, err := stream.Recv()
		if err != nil {
			fmt.Println("Subscriber has disconnected")
			s.StreamMgrCmd(newStream, "Remove")
		}
	}()

	// TODO: Remove the bs stalling and have a correct ending of the function if even necessary.
	for {
		time.Sleep(time.Second * 20)
	}

	return nil
}
