package northboundInterface

import (
	"context"
	"time"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Get(ctx context.Context, req *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("Denied a Get request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}

	log.Info("Allowed a Get request")

	notifications := make([]*gnmi.Notification, 1)
	prefix := req.GetPrefix()
	ts := time.Now().UnixNano()

	notifications[0] = &gnmi.Notification{
		Timestamp: ts,
		Prefix:    prefix,
	}

	resp := &gnmi.GetResponse{Notification: notifications}

	return resp, nil
}
