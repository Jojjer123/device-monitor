package adminInterface

import (
	"fmt"
	"net"
	"sync"

	adminServer "github.com/onosproject/device-monitor/pkg/admin/generated"
	"google.golang.org/grpc"
)

func startServer(serverWaitGroup *sync.WaitGroup) {
	defer serverWaitGroup.Done()
	// Default rpc port range is 1024-5000.
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		fmt.Println("Could not start to listen on port 4040")
	}

	s := adminServer.Server{}

	grpcServer := grpc.NewServer()

	adminServer.RegisterMonitorAdminInterfaceServer(grpcServer, &s)

	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Println("Failed to serve gRPC server on port 4040")
	}
}
