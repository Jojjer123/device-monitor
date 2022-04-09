package types

import (
	"sync"

	"github.com/openconfig/gnmi/proto/gnmi"
)

type AdminChannelMessage struct {
	RegisterFunction func(chan string, *sync.WaitGroup)
	ExecuteSetCmd    func(string, string, int) string
	Message          string
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
	Target         string
	Adapter        Adapter
	ManagerChannel <-chan string
}
