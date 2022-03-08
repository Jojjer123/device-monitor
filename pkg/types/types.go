package types

import "sync"

type AdminChannelMessage struct {
	RegisterFunction func(chan string, *sync.WaitGroup)
	Message          string
}

type Config struct {
	DevicesWithMonitoring []struct {
		DeviceName     string `yaml:"device_name"`
		DeviceIP       string `yaml:"device_ip"`
		Protocol       string `yaml:"protocol"`
		DeviceCounters []struct {
			Name     string `yaml:"name"`
			Interval int    `yaml:"interval"`
			Path     string `yaml:"path"`
		} `yaml:"device_counters"`
	} `yaml:"devices_with_monitoring"`
}
