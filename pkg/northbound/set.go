package northboundInterface

import (
	"fmt"
	"strconv"
	"time"

	"github.com/onosproject/monitor-service/pkg/configManager"
	"github.com/onosproject/monitor-service/pkg/logger"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		logger.Infof("Denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}

	fmt.Printf("Set request start: %v\n", time.Now().UnixNano())

	logger.Info("Allowed a Set request")

	var updateResult []*gnmi.UpdateResult

	// TODO: Add logging to let the user know that monitoring has started.

	// TODO: Change to 3 cases, Start, Update, and Stop. If the setRequest and deleteRequest functions are merged, there is no need for a switch.
	for _, update := range req.Update {
		if update.Path.Elem[0].Name == "Action" {
			switch update.Path.Elem[0].Key["Action"] {
			case "Start":
				updateResult = append(updateResult, s.setRequest(update.Path))
			case "Update":
				updateResult = append(updateResult, s.setRequest(update.Path))
			case "Stop":
				updateResult = append(updateResult, s.deleteRequest(update.Path))
			default:
				logger.Error("Action not found!")
			}
		}
	}

	response := gnmi.SetResponse{
		Response:  updateResult,
		Timestamp: time.Now().UnixNano(),
	}

	return &response, nil
}

// TODO: Merge with setRequest, they are similar enough to be used as one function, it only requires an extra if-statement most likely.
func (s *server) deleteRequest(path *gnmi.Path) *gnmi.UpdateResult {
	update := gnmi.UpdateResult{
		Path: &gnmi.Path{
			Element: []string{configManager.ExecuteAdminSetCmd(path.Elem[0].Key["Action"], path.Target)},
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
		logger.Error("Failed to convert ConfigIndex from string to int")
	} else {
		update := gnmi.UpdateResult{
			Path: &gnmi.Path{
				Element: []string{configManager.ExecuteAdminSetCmd(action, path.Target, configIndex)},
				Target:  path.Target,
			},
		}

		return &update
	}

	return nil
}
