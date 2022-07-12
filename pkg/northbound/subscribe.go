package northboundInterface

import (
	"time"

	"github.com/google/gnxi/utils/credentials"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/openconfig/gnmi/proto/gnmi"

	"github.com/onosproject/monitor-service/pkg/subscriptionManager"
	"github.com/onosproject/monitor-service/pkg/types"
)

func (s *server) Subscribe(stream pb.GNMI_SubscribeServer) error {
	msg, ok := credentials.AuthorizeUser(stream.Context())
	if !ok {
		log.Infof("Denied a Subscribe request: %v", msg)

		return status.Error(codes.PermissionDenied, msg)
	}

	log.Infof("Allowed a Subscribe request: %v", msg)

	subRequest, err := stream.Recv()
	if err != nil {
		log.Errorf("Failed to receive from stream: %v", err)
	}

	// fmt.Println(subRequest.GetSubscribe().Subscription[0].Path)

	for _, sub := range subRequest.GetSubscribe().Subscription {
		newStream := types.Stream{
			StreamHandle: stream,
			Target:       sub.Path.Elem,
		}

		// s.StreamMgrCmd(newStream, "Add")
		subscriptionManager.SubscriptionMgrCmd(newStream, "Add")

		go func() {
			_, err := stream.Recv()
			if err != nil {
				// fmt.Println("Subscriber has disconnected")
				log.Info("Subscriber has disconnected")
				// s.StreamMgrCmd(newStream, "Remove")
				subscriptionManager.SubscriptionMgrCmd(newStream, "Remove")
			}
		}()
	}

	// TODO: Remove the bs stalling and have a correct ending of the function if even necessary.
	for {
		time.Sleep(time.Second * 20)
	}

	return nil
}
