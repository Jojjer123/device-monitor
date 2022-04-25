package types

import (
	"github.com/golang/protobuf/proto"
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

type AdapterResponse struct {
	Entries   []SchemaEntry `protobuf:"bytes,2,opt,name=Entries,proto3"`
	Timestamp int64         `protobuf:"fixed64,1,opt,name=Timestamp,proto3"`
}

func (m *AdapterResponse) Reset()         { *m = AdapterResponse{} }
func (m *AdapterResponse) String() string { return proto.CompactTextString(m) }
func (m *AdapterResponse) ProtoMessage()  {}

type SchemaEntry struct {
	Name      string `protobuf:"bytes,2,req,name=Name,proto3"`
	Tag       string `protobuf:"bytes,2,opt,name=Tag,proto3"`
	Namespace string `protobuf:"bytes,2,opt,name=Namespace,proto3"`
	Value     string `protobuf:"bytes,2,opt,name=Value,proto3"`
}

type GatewayData struct {
	Data             string `json:"data"`
	MonitorTimestamp int64  `json:"monitorTimestamp"`
	AdapterTimestamp int64  `json:"adapterTimestamp"`
}
