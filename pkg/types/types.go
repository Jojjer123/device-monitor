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

// CONFIG STRUCTURE FOR MONITORING:
type DeviceCounter struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type IntervalCounters struct {
	Interval int             `yaml:"interval"`
	Counters []DeviceCounter `yaml:"counters"`
}

type Conf struct {
	Counters []IntervalCounters `yaml:"config"`
}

type ConfigObject struct {
	DeviceIP   string `yaml:"device_ip"`
	DeviceName string `yaml:"device_name"`
	Protocol   string `yaml:"protocol"`
	Configs    []Conf `yaml:"configs"`
}

type Counter struct {
	Name string
	Path []*gnmi.PathElem
}

type Request struct {
	Interval    int
	Counters    []Counter
	GnmiRequest *gnmi.GetRequest
}

// type Request struct {
// 	Name     string
// 	Interval int
// 	Path     []*gnmi.PathElem
// }

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
	Entries   []SchemaEntry `protobuf:"bytes,1,opt,name=Entries"`
	Timestamp int64         `protobuf:"fixed64,2,opt,name=Timestamp"`
}

func (m *AdapterResponse) Reset()         { *m = AdapterResponse{} }
func (m *AdapterResponse) String() string { return proto.CompactTextString(m) }
func (m *AdapterResponse) ProtoMessage()  {}

type SchemaEntry struct {
	Name      string `protobuf:"bytes,1,req,name=Name"`
	Tag       string `protobuf:"bytes,2,opt,name=Tag"`
	Namespace string `protobuf:"bytes,3,opt,name=Namespace"`
	Value     string `protobuf:"bytes,4,opt,name=Value"`
}

type GatewayData struct {
	Data             []byte `json:"data"`
	MonitorTimestamp int64  `json:"monitorTimestamp"`
	AdapterTimestamp int64  `json:"adapterTimestamp"`
}
