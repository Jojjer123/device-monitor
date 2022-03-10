package northboundInterface

import (
	"fmt"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		fmt.Print("Denied a Set request: ")
		fmt.Println(msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	fmt.Print("Allowed a Set request: ")
	fmt.Println(msg)

	setResponse, err := s.Server.Set(ctx, req)
	return setResponse, err
}
