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

	"github.com/onosproject/monitor-service/pkg/types"
)

func startServer(address string, executeSetCmd func(string, string, ...int) string, streamMgrCmd func(types.Stream, string) string) {
	model := gnmi.NewModel(modeldata.ModelData,
		reflect.TypeOf((*gostruct.Device)(nil)),
		gostruct.SchemaTree["Device"],
		gostruct.Unmarshal,
		gostruct.ΛEnum)

	// TODO Add credentials

	// flag.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "Supported models:\n")
	// 	for _, m := range model.SupportedModels() {
	// 		fmt.Fprintf(os.Stderr, "  %s\n", m)
	// 	}
	// 	fmt.Fprintf(os.Stderr, "\n")
	// 	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	// 	flag.PrintDefaults()
	// }

	// flag.Parse()

	// opts := credentials.ServerCredentials()
	g := grpc.NewServer( /*opts...*/ )

	var configData []byte

	var err error
	configData, err = ioutil.ReadFile("./target_configs/typical_ofsw_config.json") //*configFile)
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
