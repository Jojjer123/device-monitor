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

	confInterface "github.com/onosproject/device-monitor/pkg/config"
)

func (s *server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		fmt.Print("Denied a Set request: ")
		fmt.Println(msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	fmt.Println("Allowed a Set request")

	var updateResult []*gnmi.UpdateResult

	for _, update := range req.Update {
		if update.Path.Elem[0].Name == "Action" {
			switch update.Path.Elem[0].Key["Action"] {
			case "Create":
				// Change input param to update or path
				updateResult = append(updateResult, s.setRequest(update.Path))
			case "Update":
				updateResult = append(updateResult, s.setRequest(update.Path))
			case "Change config":
				updateResult = append(updateResult, s.updateConfigRequest(req))
			case "Delete":
				updateResult = append(updateResult, s.deleteRequest(update.Path))
			default:
				fmt.Println("Action not found!")
			}
		}
	}

	response := gnmi.SetResponse{
		Response:  updateResult,
		Timestamp: time.Now().UnixNano(),
	}

	return &response, nil
}

func (s *server) deleteRequest(path *gnmi.Path) *gnmi.UpdateResult {
	update := gnmi.UpdateResult{
		Path: &gnmi.Path{
			Element: []string{s.ExecuteSetCmd(path.Elem[0].Key["Action"], path.Target)},
			Target:  path.Target,
		},
	}

	return &update
}

func (s *server) setRequest(path *gnmi.Path) *gnmi.UpdateResult {
	var configIndex int
	var err error

	action := path.Elem[0].Key["Action"]

	if path.Elem[1].Name == "ConfigIndex" {
		configIndex, err = strconv.Atoi(path.Elem[1].Key["ConfigIndex"])
	}

	if err != nil {
		fmt.Println("Failed to convert ConfigIndex from string to int")
	} else {
		update := gnmi.UpdateResult{
			Path: &gnmi.Path{
				Element: []string{s.ExecuteSetCmd(action, path.Target, configIndex)},
				Target:  path.Target,
			},
		}

		return &update
	}

	return nil
}

func (s *server) updateConfigRequest(req *gnmi.SetRequest) *gnmi.UpdateResult {
	err := confInterface.UpdateConfig(req)
	if err != nil {
		fmt.Println("Failed to update configuration!")
	}

	// TODO: Add correct result

	return &gnmi.UpdateResult{}
}
