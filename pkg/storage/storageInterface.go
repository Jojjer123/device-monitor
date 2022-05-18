package config

import (
	"context"
	"encoding/json"
	"fmt"

	// "sync"
	"time"

	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"gopkg.in/yaml.v2"

	types "github.com/onosproject/monitor-service/pkg/types"
)

const (
	GetRequest_MONITOR_CONFIG  gnmi.GetRequest_DataType = 5 // Configuration for a switch.
	GetRequest_MONITOR_ADAPTER gnmi.GetRequest_DataType = 6 // Adapter for a protocol.
)

func GetConfig(target string) types.ConfigObject {
	ctx := context.Background()

	// Create a gNMI client, if credentials is required, implement it here. Storage-service does not
	// offer secure communication yet. Don't forget to change port if changing to secure communication.
	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{"storage-service:11161"},
		Target:      "storage-service",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	r := &gnmi.GetRequest{
		Type: GetRequest_MONITOR_CONFIG,
		Path: []*gnmi.Path{
			{
				Target: target,
			},
		},
	}

	response, err := c.(*gclient.Client).Get(ctx, r)

	if err != nil {
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

	// Create a gNMI client, if secure connection is required, implement it here. Storage-service does not
	// offer secure communication yet. Don't forget to change port if changing to secure communication.
	c, err := gclient.New(ctx, client.Destination{
		Addrs:       []string{"storage-service:11161"},
		Target:      "storage-service",
		Timeout:     time.Second * 5,
		Credentials: nil,
		TLS:         nil,
	})

	if err != nil {
		fmt.Print("Could not create a gNMI client: ")
		fmt.Println(err)
	}

	r := &gnmi.GetRequest{
		Type: GetRequest_MONITOR_ADAPTER,
		Path: []*gnmi.Path{
			{
				Target: protocol,
			},
		},
	}

	response, err := c.(*gclient.Client).Get(ctx, r)

	if err != nil {
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
