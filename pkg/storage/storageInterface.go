package config

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"gopkg.in/yaml.v2"

	types "github.com/onosproject/monitor-service/pkg/types"
)

// TODO: Remove all bs init functions that doesn't do shit.
func ConfigInterface(waitGroup *sync.WaitGroup) {
	// fmt.Println("ConfigInterface started")
	defer waitGroup.Done()

	// var c client.Impl
	// var err error

	// fmt.Println(response)
}

// TODO: Implement cert usage

func GetConfig(target string) types.ConfigObject {
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

	r := &gnmi.GetRequest{
		Type: 5,
		Path: []*gnmi.Path{
			{
				Target: target,
			},
		},
	}

	response, err := c.(*gclient.Client).Get(ctx, r)

	if err != nil {
		// fmt.Errorf("target returned RPC error for Get(%q): %v", r.String(), err)
		fmt.Print("Target returned RPC error for Get(")
		fmt.Print(r.String())
		fmt.Print("): ")
		fmt.Println(err)
	}

	var config types.ConfigObject
	yaml.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &config)

	return config
}

func GetAdapter(protocol string) types.Adapter {
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

	r := &gnmi.GetRequest{
		Type: 6,
		Path: []*gnmi.Path{
			{
				Target: protocol,
			},
		},
	}

	response, err := c.(*gclient.Client).Get(ctx, r)

	if err != nil {
		// fmt.Errorf("target returned RPC error for Get(%q): %v", r.String(), err)
		fmt.Print("Target returned RPC error for Get(")
		fmt.Print(r.String())
		fmt.Print("): ")
		fmt.Println(err)
	}

	var adapter types.Adapter
	err = json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &adapter)
	if err != nil {
		fmt.Print("Failed to unmarshal adapter")
		fmt.Println(err)
	}

	return adapter
}

func UpdateConfig(req *gnmi.SetRequest) error {
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

	// r := &gpb.Set{}

	// err = proto.UnmarshalText(`path: <target: '`+target+`', elem: <name: 'interfaces'>
	// 				elem: <name: 'interface' key: <key: 'name' value: 'sw0p5'>>>`, r)

	// if err != nil {
	// 	// fmt.Errorf("unable to parse gnmi.GetRequest: %v", err)
	// 	fmt.Print("Unable to parse gnmi.GetRequest: ")
	// 	fmt.Println(err)
	// }

	response, err := c.(*gclient.Client).Set(ctx, req)

	if err != nil {
		// fmt.Errorf("target returned RPC error for Get(%q): %v", r.String(), err)
		fmt.Print("Target returned RPC error for Get(")
		fmt.Print(req.String())
		fmt.Print("): ")
		fmt.Println(err)
	}

	// var config types.ConfigRequest
	// yaml.Unmarshal(response.Notification[0].Update[0].Value.Value, &config)

	fmt.Println(response)

	return err
}
