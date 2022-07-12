package Storage

import (
	"fmt"

	"github.com/onosproject/monitor-service/pkg/logger"
	"google.golang.org/protobuf/proto"

	// types "github.com/onosproject/monitor-service/pkg/types"

	"github.com/onosproject/monitor-service/pkg/proto/adapter"
	conf "github.com/onosproject/monitor-service/pkg/proto/monitor-config"
)

// const (
// 	GetRequest_MONITOR_CONFIG  gnmi.GetRequest_DataType = 5 // Configuration for a switch.
// 	GetRequest_MONITOR_ADAPTER gnmi.GetRequest_DataType = 6 // Adapter for a protocol.
// )

var log = logger.GetLogger()

func GetConfig(target string) (*conf.MonitorConfig, error) {
	urn := fmt.Sprintf("configurations.monitor-config.%s", target)

	// log.Infof("Get config for %v, from ext service: %v\n", target, time.Now().UnixNano())

	// log.Infof("Getting config using urn: %v", urn)

	rawConf, err := getRawDataFromStore(urn)
	if err != nil {
		log.Errorf("Failed getting config from store: %v", err)
		return &conf.MonitorConfig{}, err
	}

	// log.Infof("Got raw conf from store: %v", rawConf)

	// log.Infof("Received config for %v, from ext service: %v\n", target, time.Now().UnixNano())

	var config = &conf.Config{}

	if err = proto.Unmarshal(rawConf, config); err != nil {
		log.Errorf("Failed unmarshaling config from store: %v", err)
		return &conf.MonitorConfig{}, err
	}

	log.Infof("Config umarshaled into: %v", config)

	return config.Devices[0], nil
}

func GetAdapter(protocol string) (*adapter.Adapter, error) {
	urn := fmt.Sprintf("configurations.adapter.%s", protocol)

	// log.Infof("Get adapter for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

	rawAdapter, err := getRawDataFromStore(urn)
	if err != nil {
		log.Errorf("Failed getting adapter from store: %v", err)
		return &adapter.Adapter{}, err
	}

	// log.Infof("Received adapter for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

	var adapterRef = &adapter.Adapter{}

	if err = proto.Unmarshal(rawAdapter, adapterRef); err != nil {
		log.Errorf("Failed unmarshaling adapter from store: %v", err)
		return &adapter.Adapter{}, err
	}

	return adapterRef, nil
}

// func GetConfig(target string) types.ConfigObject {
// 	ctx := context.Background()

// 	// Create a gNMI client, if credentials is required, implement it here. Storage-service does not
// 	// offer secure communication yet. Don't forget to change port if changing to secure communication.
// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       []string{"storage-service:11161"},
// 		Target:      "storage-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		log.Errorf("Could not create a gNMI client: %v", err)
// 	}

// 	r := &gnmi.GetRequest{
// 		Type: GetRequest_MONITOR_CONFIG,
// 		Path: []*gnmi.Path{
// 			{
// 				Target: target,
// 			},
// 		},
// 	}

// 	log.Infof("Get config for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

// 	response, err := c.(*gclient.Client).Get(ctx, r)

// 	log.Infof("Received config for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

// 	if err != nil {
// 		log.Errorf("Target returned RPC error for Get(%v): %v", r.String(), err)
// 	}

// 	var config types.ConfigObject
// 	err = yaml.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &config)
// 	if err != nil {
// 		log.Errorf("Could not unmarshal config: %v", err)
// 	}

// 	return config
// }

// func GetAdapter(protocol string) types.Adapter {
// 	ctx := context.Background()

// 	// Create a gNMI client, if secure connection is required, implement it here. Storage-service does not
// 	// offer secure communication yet. Don't forget to change port if changing to secure communication.
// 	c, err := gclient.New(ctx, client.Destination{
// 		Addrs:       []string{"storage-service:11161"},
// 		Target:      "storage-service",
// 		Timeout:     time.Second * 5,
// 		Credentials: nil,
// 		TLS:         nil,
// 	})

// 	if err != nil {
// 		log.Errorf("Could not create a gNMI client: %v", err)
// 	}

// 	r := &gnmi.GetRequest{
// 		Type: GetRequest_MONITOR_ADAPTER,
// 		Path: []*gnmi.Path{
// 			{
// 				Target: protocol,
// 			},
// 		},
// 	}

// 	log.Infof("Get adapter for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

// 	response, err := c.(*gclient.Client).Get(ctx, r)

// 	log.Infof("Received adapter for %v, from ext service: %v\n", r.Path[0].Target, time.Now().UnixNano())

// 	if err != nil {
// 		log.Errorf("Target returned RPC error for Get(%v): %v", r.String(), err)
// 	}

// 	var adapter types.Adapter
// 	err = json.Unmarshal(response.Notification[0].Update[0].Val.GetBytesVal(), &adapter)
// 	if err != nil {
// 		log.Errorf("Failed to unmarshal adapter: %v", err)
// 	}

// 	return adapter
// }
