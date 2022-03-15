package northboundInterface

import (
	"fmt"
	"time"

	// "github.com/docker/engine/api/types/time"
	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/openconfig/gnmi/value"
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
	fmt.Println("Allowed a Set request")

	cmd, err := value.ToScalar(req.Update[0].Val)

	if err != nil {
		fmt.Println("Failed to convert gnmi.TypedValue to scalar")
	}

	var updateResult []*gnmi.UpdateResult

	update := gnmi.UpdateResult{
		Path: &gnmi.Path{
			Element: []string{s.ExecuteSetCmd(cmd.(string))},
			Target:  req.Update[0].Path.Target,
		},
	}

	updateResult = append(updateResult, &update)

	response := gnmi.SetResponse{
		Response:  updateResult,
		Timestamp: time.Now().UnixNano(),
	}

	return &response, nil
	// setResponse, err := s.Server.Set(ctx, req)
	// return setResponse, err
}
