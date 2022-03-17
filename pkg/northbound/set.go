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
				{
					updateResult = append(updateResult, s.createRequest(req))
				}
			case "Update":
				{
					updateResult = append(updateResult, s.updateDeviceMonitorRequest(req))
				}
			case "Change config":
				{
					updateResult = append(updateResult, s.updateConfigRequest(req))
				}
			default:
				{
					fmt.Println("Action not found!")
				}
			}
		}
	}

	response := gnmi.SetResponse{
		Response:  updateResult,
		Timestamp: time.Now().UnixNano(),
	}

	return &response, nil
}

func (s *server) createRequest(req *gnmi.SetRequest) *gnmi.UpdateResult {
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
				Element: []string{s.ExecuteSetCmd(action, req.Update[0].Path.Target, configIndex)},
				Target:  req.Update[0].Path.Target,
			},
		}
	}

	return &update
}

func (s *server) updateConfigRequest(req *gnmi.SetRequest) *gnmi.UpdateResult {
	err := confInterface.UpdateConfig(req)
	if err != nil {
		fmt.Println("Failed to update configuration!")
	}

	// TODO: Add correct result

	return &gnmi.UpdateResult{}
}

func (s *server) updateDeviceMonitorRequest(req *gnmi.SetRequest) *gnmi.UpdateResult {

	return &gnmi.UpdateResult{}
}
