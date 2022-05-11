package requestBuilder

import (
	"strings"
	"sync"

	// "github.com/google/gnxi/gnmi"
	storageInterface "github.com/onosproject/monitor-service/pkg/storage"
	types "github.com/onosproject/monitor-service/pkg/types"
	"github.com/openconfig/gnmi/proto/gnmi"
)

// TODO: Remove this bs init function that is doing nothing.
func RequestBuilder(waitGroup *sync.WaitGroup) {
	// fmt.Println("RequestBuilder started")
	defer waitGroup.Done()

}

func GetConfig(target string, configSelected int) ([]types.Request, types.Adapter) {
	conf := storageInterface.GetConfig(target)

	// TODO: Add check for empty config, and dont crash if that is the case.

	var requests []types.Request

	// fmt.Println("----CONFIG----")
	// fmt.Printf("%v\n", conf.Configs[configSelected])
	// fmt.Println("--------------")

	// TODO: Change from single reqeustObj to batchObj that is based on interval
	for _, intCounters := range conf.Configs[configSelected].Counters {
		request := types.Request{
			Interval: intCounters.Interval,
		}

		for _, counter := range intCounters.Counters {
			request.Counters = append(request.Counters, types.Counter{
				Name: counter.Name,
				Path: getPathFromString(counter.Path),
			})
		}
		// requestObj := types.Request{
		// 	Name:     req.Name,
		// 	Interval: req.Interval,
		// 	Path:     getPathFromString(req.Path),
		// }

		requests = append(requests, request)
	}

	var adapter types.Adapter

	if conf.Protocol != "GNMI" {
		adapter = storageInterface.GetAdapter(conf.Protocol)
	}

	return requests, adapter
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
