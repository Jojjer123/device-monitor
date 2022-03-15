package requestBuilder

import (
	"fmt"
	"sync"

	confInterface "github.com/onosproject/device-monitor/pkg/config"
	types "github.com/onosproject/device-monitor/pkg/types"
	"gopkg.in/yaml.v2"
)

func RequestBuilder(waitGroup *sync.WaitGroup) {
	fmt.Println("RequestBuilder started")
	defer waitGroup.Done()

}

func GetRequest(target string, confType string) []types.Request {
	conf := confInterface.GetConfig(target)

	var requests []types.Request

	switch confType {
	case "default":
		{
			for _, req := range conf.DefaultConfig.DeviceCounters {
				var requestObj types.Request
				reqBytes, err := yaml.Marshal(req)
				if err != nil {
					fmt.Println("Failed to convert device counter to byte slice!")
				}
				err = yaml.Unmarshal(reqBytes, &requestObj)
				if err != nil {
					fmt.Println("Failed to convert byte slice to request struct")
				}
				requests = append(requests, requestObj)
			}
		}
	case "temporary":
		{

		}
	default:
		{
			fmt.Println("Configuration type is undefined!")
		}
	}

	return requests
}

// func getConfig(target string) types.Config {

// 	return types.Config{}
// }
