package types

import (
	"github.com/openconfig/gnmi/proto/gnmi"
)

type ConfigAdminChannelMessage struct {
	// RegisterFunction func(chan string, *sync.WaitGroup)
	ExecuteSetCmd func(string, string, ...int) string
	// Message       string
}

type StreamMgrChannelMessage struct {
	ManageCmd func(Stream, string) string
	// ExecuteSetCmd    func(string, string, ...int) string
	// Message          string
}

type Stream struct {
	StreamHandle gnmi.GNMI_SubscribeServer
	Target       []*gnmi.PathElem
}

type ConfigRequest struct {
	DeviceIP   string `yaml:"device_ip"`
	DeviceName string `yaml:"device_name"`
	Protocol   string `yaml:"protocol"`
	Configs    []struct {
		DeviceCounters []struct {
			Name     string `yaml:"name"`
			Interval int    `yaml:"interval"`
			Path     string `yaml:"path"`
		} `yaml:"device_counters"`
	} `yaml:"configs"`
}

type Request struct {
	Name     string
	Interval int
	Path     []*gnmi.PathElem
}

type Adapter struct {
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
}

type DeviceMonitor struct {
	Target          string
	Adapter         Adapter
	Requests        []Request
	RequestsChannel chan []Request
	ManagerChannel  chan string
}

// The following types are used for deconstructing data from the adapter.

type SchemaTree struct {
	Name      string
	Namespace string
	Children  []*SchemaTree
	Parent    *SchemaTree
	Value     string
}

type Schema struct {
	Entries []SchemaEntry
}

type SchemaEntry struct {
	Name      string
	Tag       string
	Namespace string
	Value     string
}

// type GatewayData struct {
// 	Data      string `json:"data"`
// 	Timestamp int64  `json:"timestamp"`
// }
