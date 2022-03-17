package northboundInterface

import (
	"fmt"
	"strconv"
	"time"

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
	fmt.Println("Allowed a Set request")

	// cmd, err := value.ToScalar(req.Update[0].Val)

	// if err != nil {
	// 	fmt.Println("Failed to convert gnmi.TypedValue to scalar")
	// }

	var updateResult []*gnmi.UpdateResult

	var action string
	var configIndex int
	var err error

	if req.Update[0].Path.Elem[0].Name == "Action" {
		action = req.Update[0].Path.Elem[0].Key["Action"]
	}
	if req.Update[0].Path.Elem[1].Name == "ConfigIndex" {
		configIndex, err = strconv.Atoi(req.Update[0].Path.Elem[1].Key["ConfigIndex"])

		if err != nil {
			fmt.Println("Failed to convert ConfigIndex from string to int")
		}
	}

	var update gnmi.UpdateResult

	if configIndex > len(req.Update[0].Path.Elem)-1 || configIndex < 0 {
		fmt.Println("Configuration index is out of bounds!")
	} else {
		update = gnmi.UpdateResult{
			Path: &gnmi.Path{
				Element: []string{s.ExecuteSetCmd(action /*cmd.(string)*/, req.Update[0].Path.Target, configIndex /*index of selected config*/)},
				Target:  req.Update[0].Path.Target,
			},
		}
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
