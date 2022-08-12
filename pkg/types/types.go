package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/onosproject/monitor-service/pkg/proto/adapter"
	"github.com/openconfig/gnmi/proto/gnmi"
)

type ConfigAdminChannelMessage struct {
	ExecuteSetCmd func(string, string, ...int) string
}

type StreamMgrChannelMessage struct {
	ManageCmd func(Stream, string) string
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

// DEVICE MONITOR
type Counter struct {
	Name string
	Path []*gnmi.PathElem
}

type Request struct {
	Interval    int
	Counters    []Counter
	GnmiRequest *gnmi.GetRequest
}

type DeviceMonitor struct {
	DeviceName      string
	Target          string
	Adapter         *adapter.Adapter
	Requests        []Request
	RequestsChannel chan []Request
	ManagerChannel  chan string
}

// ALL OF THE BELOW STRUCTURES SHOULD BE IMPLEMENTED THROUGH PROTOBUF
// OR THEY SHOULD BE REMOVED WITH A NATIVE GNMI IMPLEMENTATION
// The following types are used for deconstructing data from the adapter
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

type Dictionary map[string]interface{}

type GatewayData struct {
	Data             []Dictionary `json:"data"`
	MonitorTimestamp int64        `json:"monitorTimestamp"`
	AdapterTimestamp int64        `json:"adapterTimestamp"`
}
