package northboundInterface

import (
	"github.com/google/gnxi/gnmi"
	"github.com/openconfig/ygot/ygot"

	"github.com/onosproject/monitor-service/pkg/types"
)

type server struct {
	*gnmi.Server
	Model        *gnmi.Model
	configStruct ygot.ValidatedGoStruct
	StreamMgrCmd func(types.Stream, string) string
}

func newServer(model *gnmi.Model, config []byte) (*server, error) {
	s, err := gnmi.NewServer(model, config, nil)

	if err != nil {
		return nil, err
	}

	newconfig, _ := model.NewConfigStruct(config)

	server := server{
		Server:       s,
		Model:        model,
		configStruct: newconfig,
	}

	return &server, nil
}
