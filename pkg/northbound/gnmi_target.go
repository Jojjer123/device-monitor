package northboundInterface

import (
	"io/ioutil"
	"net"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/google/gnxi/gnmi"
	"github.com/google/gnxi/gnmi/modeldata"
	"github.com/google/gnxi/gnmi/modeldata/gostruct"

	pb "github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/credentials"
)

func startServer(secure bool, address string) {
	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*gostruct.Device)(nil)),
		gostruct.SchemaTree["Device"],
		gostruct.Unmarshal,
		gostruct.ΛEnum)

	var g *grpc.Server

	// Create server with credentials, they are COPIED from gnxi-simulators, so they SHOULD be replaced
	if secure {
		creds, err := credentials.NewServerTLSFromFile("certs/localhost.crt", "certs/localhost.key")
		if err != nil {
			log.Errorf("Failed to load credentials: %v\n", err)
		}

		g = grpc.NewServer(grpc.Creds(creds))
	} else {
		g = grpc.NewServer()
	}

	configData, err := ioutil.ReadFile("./target_configs/typical_ofsw_config.json")
	if err != nil {
		log.Errorf("Error in reading config file: %v", err)
	}

	s, err := newServer(model, configData)

	if err != nil {
		log.Errorf("Error in creating gnmi target: %v", err)
	}
	pb.RegisterGNMIServer(g, s)
	reflection.Register(g)

	log.Infof("Starting gNMI agent to listen on %v", address)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("Failed to listen: %v", err)
	}

	log.Infof("Starting gNMI agent to serve on %v", address)
	if err := g.Serve(listen); err != nil {
		log.Errorf("Failed to serve: %v", err)
	}
}
