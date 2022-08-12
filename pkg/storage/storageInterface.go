package Storage

import (
	"fmt"

	"github.com/onosproject/monitor-service/pkg/logger"
	"google.golang.org/protobuf/proto"

	"github.com/onosproject/monitor-service/pkg/proto/adapter"
	conf "github.com/onosproject/monitor-service/pkg/proto/monitor-config"
)

var log = logger.GetLogger()

// Gets monitor configuration from k/v store
func GetConfig(target string) (*conf.MonitorConfig, error) {
	urn := fmt.Sprintf("configurations.monitor-config.%s", target)

	// Get raw data from k/v store
	rawConf, err := getRawDataFromStore(urn)
	if err != nil {
		log.Errorf("Failed getting config from store: %v", err)
		return &conf.MonitorConfig{}, err
	}

	var config = &conf.Config{}

	// Deserialize raw data to config
	if err = proto.Unmarshal(rawConf, config); err != nil {
		log.Errorf("Failed unmarshaling config from store: %v", err)
		return &conf.MonitorConfig{}, err
	}

	return config.Devices[0], nil
}

// Gets adapter from k/v store
func GetAdapter(protocol string) (*adapter.Adapter, error) {
	urn := fmt.Sprintf("configurations.adapter.%s", protocol)

	// Get raw data from k/v store
	rawAdapter, err := getRawDataFromStore(urn)
	if err != nil {
		log.Errorf("Failed getting adapter from store: %v", err)
		return &adapter.Adapter{}, err
	}

	var adapterRef = &adapter.Adapter{}

	// Deserialize raw data to adapter
	if err = proto.Unmarshal(rawAdapter, adapterRef); err != nil {
		log.Errorf("Failed unmarshaling adapter from store: %v", err)
		return &adapter.Adapter{}, err
	}

	return adapterRef, nil
}
