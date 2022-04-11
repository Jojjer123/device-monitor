package requestBuilder

import (
	"fmt"
	"strings"
	"sync"

	// "github.com/google/gnxi/gnmi"
	confInterface "github.com/onosproject/device-monitor/pkg/config"
	types "github.com/onosproject/device-monitor/pkg/types"
	"github.com/openconfig/gnmi/proto/gnmi"
)

func RequestBuilder(waitGroup *sync.WaitGroup) {
	fmt.Println("RequestBuilder started")
	defer waitGroup.Done()

}

func GetConfig(target string, configSelected int) ([]types.Request, types.Adapter, string) {
	conf := confInterface.GetConfig(target)

	// TODO: Add check for empty config, and dont crash if that is the case.

	var requests []types.Request

	for _, req := range conf.Configs[configSelected].DeviceCounters {
		requestObj := types.Request{
			Name:     req.Name,
			Interval: req.Interval,
			Path:     getPathFromString(req.Path),
		}
		// reqBytes, err := yaml.Marshal(req)
		// if err != nil {
		// 	fmt.Println("Failed to convert device counter to byte slice!")
		// }
		// err = yaml.Unmarshal(reqBytes, &requestObj)
		// if err != nil {
		// 	fmt.Println("Failed to convert byte slice to request struct")
		// }
		requests = append(requests, requestObj)
	}

	var adapter types.Adapter

	if conf.Protocol != "GNMI" {
		adapter = confInterface.GetAdapter(conf.Protocol)
	}

	return requests, adapter, conf.DeviceIP
}

//  <name: 'interfaces' key: <key: 'namespace' value: 'urn:ietf:params:xml:ns:yang:ietf-interfaces'>>
//  <name: 'interface'>
//  <name: 'sw0p1'>
//  <name: 'ethernet' <key: 'name' value: 'urn:ieee:std:802.3:yang:ieee802-ethernet-interface'>>
//  <name: 'statistics'>
//  <name: 'frame'>
//  <name: 'in-total-frames'>
func getPathFromString(path string) []*gnmi.PathElem {
	if !strings.Contains(path, "elem:") {
		return nil
	}

	var pathElems []*gnmi.PathElem
	for index, elem := range strings.Split(path, "elem:") {
		if index == 0 {
			continue
		}

		// fmt.Println(elem)
		// fmt.Println("--------------")
		tok := strings.Split(elem, "'")

		// fmt.Println(tok)
		// fmt.Println(tok[1])

		newElem := &gnmi.PathElem{
			Name: tok[1],
		}

		// Contains key.
		if len(tok) > 3 {
			// fmt.Println(tok[3])
			// fmt.Println(tok[5])
			keyMap := make(map[string]string)
			keyMap[tok[3]] = tok[5]
			newElem.Key = keyMap
		}

		pathElems = append(pathElems, newElem)
	}

	return pathElems
}

// func getConfig(target string) types.Config {

// 	return types.Config{}
// }
