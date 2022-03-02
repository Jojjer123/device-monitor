package adminServer

import (
	"fmt"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) MonitorDevice(ctx context.Context, message *MonitorMessage) (*MonitorResponse, error) {
	fmt.Println("Received message from client: ", message.Action, " ", message.Target)
	return &MonitorResponse{Response: "Successfully created monitor"}, nil
}
