package requestBuilder

import (
	"fmt"
	"sync"
)

func RequestBuilder(waitGroup *sync.WaitGroup) {
	fmt.Println("RequestBuilder started")
	defer waitGroup.Done()
}
