package configInterface

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	gpb "github.com/openconfig/gnmi/proto/gnmi"
	"gopkg.in/yaml.v2"

	types "github.com/onosproject/device-monitor/pkg/types"
)

func ConfigInterface(waitGroup *sync.WaitGroup) {
	fmt.Println("ConfigInterface started")
	defer waitGroup.Done()

	// TODO: Implement cert usage

	// var c client.Impl
	// var err error

	ctx := context.Background()

	address := []string{"storage-service:11161"}

	c, err := gclient.New(ctx, client.Destination{
		Addrs:       address,
		Target:      "storage-service",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		// fmt.Errorf("could not create a gNMI client: %v", err)
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	// "path: <target: 'storage-service', elem: <name: 'interfaces'> elem: <name: 'interface' key: <key: 'name' value: 'sw0p5'>>>"

	r := &gpb.GetRequest{}

	if err := proto.UnmarshalText("path: <target: 'storage-service', elem: <name: 'interfaces'> elem: <name: 'interface' key: <key: 'name' value: 'sw0p5'>>>", r); err != nil {
		// fmt.Errorf("unable to parse gnmi.GetRequest: %v", err)
		fmt.Print("Unable to parse gnmi.GetRequest: ")
		fmt.Println(err)
	}

	response, err := c.(*gclient.Client).Get(ctx, r)

	if err != nil {
		// fmt.Errorf("target returned RPC error for Get(%q): %v", r.String(), err)
		fmt.Print("Target returned RPC error for Get(")
		fmt.Print(r.String())
		fmt.Print("): ")
		fmt.Println(err)
	}

	var config types.Config
	yaml.Unmarshal(response.Notification[0].Update[0].Value.Value, &config)

	fmt.Println(config.DevicesWithMonitoring[0].DeviceCounters[1].Name,
		config.DevicesWithMonitoring[0].DeviceCounters[1].Interval,
		config.DevicesWithMonitoring[0].DeviceCounters[1].Path)

	// fmt.Println(response)
}
