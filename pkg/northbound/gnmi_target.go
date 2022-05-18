package northboundInterface

import (
	"fmt"
	"io/ioutil"
	"net"
	"reflect"

	// "time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/google/gnxi/gnmi"
	"github.com/google/gnxi/gnmi/modeldata"
	"github.com/google/gnxi/gnmi/modeldata/gostruct"

	pb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/credentials"

	"github.com/onosproject/monitor-service/pkg/types"
)

func startServer(secure bool, address string, executeSetCmd func(string, string, ...int) string, streamMgrCmd func(types.Stream, string) string) {
	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*gostruct.Device)(nil)),
		gostruct.SchemaTree["Device"],
		gostruct.Unmarshal,
		gostruct.ΛEnum)

	var g *grpc.Server

	// Create server with credentials, they are COPIED from gnxi-simulators, so they SHOULD be replaced.
	if secure {
		creds, err := credentials.NewServerTLSFromFile("certs/localhost.crt", "certs/localhost.key")
		if err != nil {
			fmt.Printf("Failed to load credentials: %v\n", err)
		}

		g = grpc.NewServer(grpc.Creds(creds))
	} else {
		g = grpc.NewServer()
	}

	configData, err := ioutil.ReadFile("./target_configs/typical_ofsw_config.json") //*configFile)
	if err != nil {
		// log.Fatalf("Error in reading config file: %v", err)
		fmt.Print("Error in reading config file: ")
		fmt.Println(err)
	}

	s, err := newServer(model, configData)

	s.ExecuteSetCmd = executeSetCmd
	s.StreamMgrCmd = streamMgrCmd

	if err != nil {
		// log.Fatalf("Error in creating gnmi target: %v", err)
		fmt.Print("Error in creating gnmi target: ")
		fmt.Println(err)
	}
	pb.RegisterGNMIServer(g, s)
	reflection.Register(g)

	// log.Infof("Starting gNMI agent to listen on %s", *bindAddr)
	fmt.Print("Starting gNMI agent to listen on ")
	fmt.Println(address)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		// log.Fatalf("Failed to listen: %v", err)
		fmt.Print("Failed to listen: ")
		fmt.Println(err)
	}

	// log.Infof("Starting gNMI agent to serve on %s", *bindAddr)
	fmt.Print("Starting gNMI agent to serve on ")
	fmt.Println(address)
	if err := g.Serve(listen); err != nil {
		// log.Fatalf("Failed to serve: %v", err)
		fmt.Print("Failed to serve ")
		fmt.Println(err)
	}
}
