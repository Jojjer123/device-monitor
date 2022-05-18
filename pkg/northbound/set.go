package northboundInterface

import (
	"strconv"
	"time"

	"github.com/onosproject/monitor-service/pkg/logger"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	// storageInterface "github.com/onosproject/monitor-service/pkg/storage"
)

func (s *server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		// fmt.Print("Denied a Set request: ")
		// fmt.Println(msg)
		logger.Infof("Denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}
	logger.Info("Allowed a Set request")

	var updateResult []*gnmi.UpdateResult

	// TODO: Add logging to let the user know that monitoring has started.

	// TODO: Change to 3 cases, Start, Update, and Stop. If the setRequest and deleteRequest functions are merged, there is no need for a switch.
	for _, update := range req.Update {
		if update.Path.Elem[0].Name == "Action" {
			switch update.Path.Elem[0].Key["Action"] {
			case "Create":
				// Change input param to update or path
				updateResult = append(updateResult, s.setRequest(update.Path))
			case "Update":
				updateResult = append(updateResult, s.setRequest(update.Path))
			// case "Change config":
			// 	updateResult = append(updateResult, s.updateConfigRequest(req))
			case "Delete":
				updateResult = append(updateResult, s.deleteRequest(update.Path))
			default:
				// fmt.Println("Action not found!")
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
		// fmt.Println("Failed to convert ConfigIndex from string to int")
		logger.Error("Failed to convert ConfigIndex from string to int")
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
