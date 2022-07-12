package config

import (
	"context"
	"encoding/json"

	// "sync"
	"time"

	"github.com/openconfig/gnmi/client"
	gclient "github.com/openconfig/gnmi/client/gnmi"
	"github.com/openconfig/gnmi/proto/gnmi"
	"gopkg.in/yaml.v2"

	"github.com/onosproject/monitor-service/pkg/logger"
	types "github.com/onosproject/monitor-service/pkg/types"
)

const (
	GetRequest_MONITOR_CONFIG  gnmi.GetRequest_DataType = 5 // Configuration for a switch.
	GetRequest_MONITOR_ADAPTER gnmi.GetRequest_DataType = 6 // Adapter for a protocol.
)

var log = logger.GetLogger()

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
		log.Errorf("Could not create a gNMI client: %v", err)
	}

	r := &gnmi.GetRequest{
		Type: GetRequest_MONITOR_CONFIG,
		Path: []*gnmi.Path{
			{
				Target: target,
			},
		},
	}

	log.Infof("Get config for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

	response, err := c.(*gclient.Client).Get(ctx, r)

	log.Infof("Received config for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

	if err != nil {
		log.Errorf("Target returned RPC error for Get(%v): %v", r.String(), err)
	}

	var config types.ConfigObject
	err = yaml.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &config)
	if err != nil {
		log.Errorf("Could not unmarshal config: %v", err)
	}

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
		log.Errorf("Could not create a gNMI client: %v", err)
	}

	r := &gnmi.GetRequest{
		Type: GetRequest_MONITOR_ADAPTER,
		Path: []*gnmi.Path{
			{
				Target: protocol,
			},
		},
	}

	log.Infof("Get adapter for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

	response, err := c.(*gclient.Client).Get(ctx, r)

	log.Infof("Received adapter for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

	if err != nil {
		log.Errorf("Target returned RPC error for Get(%v): %v", r.String(), err)
	}

	var adapter types.Adapter
	err = json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &adapter)
	if err != nil {
		log.Errorf("Failed to unmarshal adapter: %v", err)
	}

	return adapter
}
