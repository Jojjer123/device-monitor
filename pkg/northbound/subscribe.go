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

	for _, sub := range subRequest.GetSubscribe().Subscription {
		newStream := types.Stream{
			StreamHandle: stream,
			Target:       sub.Path.Elem,
		}

		subscriptionManager.SubscriptionMgrCmd(newStream, "Add")

		go func() {
			_, err := stream.Recv()
			if err != nil {
				log.Info("Subscriber has disconnected")
				subscriptionManager.SubscriptionMgrCmd(newStream, "Remove")
			}
		}()
	}

	// TODO: Replace with proper methods, such as a sync.WaitGroup
	// This function should not return???
	for {
		time.Sleep(time.Second * 20)
	}
}
