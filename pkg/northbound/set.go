package northboundInterface

import (
	"strconv"
	"time"

	"github.com/onosproject/monitor-service/pkg/configManager"

	"github.com/google/gnxi/utils/credentials"
	"github.com/openconfig/gnmi/proto/gnmi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *server) Set(ctx context.Context, req *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	msg, ok := credentials.AuthorizeUser(ctx)
	if !ok {
		log.Infof("Denied a Set request: %v", msg)
		return nil, status.Error(codes.PermissionDenied, msg)
	}

	log.Infof("Allowed a Set request: %v", msg)

	var updateResult []*gnmi.UpdateResult

	// Switch over each update in the request
	for _, update := range req.Update {
		if update.Path.Elem[0].Name == "Action" {
			// TODO: Start and Update should be merged???
			switch update.Path.Elem[0].Key["Action"] {
			case "Start":
				updateResult = append(updateResult, s.setRequest(update.Path))
			case "Update":
				updateResult = append(updateResult, s.setRequest(update.Path))
			case "Stop":
				updateResult = append(updateResult, s.deleteRequest(update.Path))
			default:
				log.Error("Action not found!")
			}
		}
	}

	response := gnmi.SetResponse{
		Response:  updateResult,
		Timestamp: time.Now().UnixNano(),
	}

	return &response, nil
}

// TODO: Merge with setRequest, they are similar enough to be used as one function, it only requires an extra if-statement most likely
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
		log.Error("Failed to convert ConfigIndex from string to int")
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
