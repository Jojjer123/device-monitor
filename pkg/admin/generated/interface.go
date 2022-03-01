package adminServer

import (
	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) MonitorDevice(ctx context.Context, message *MonitorMessage) (*MonitorResponse, error) {
	return &MonitorResponse{Respone: "Successfully created monitor"}, nil
}
