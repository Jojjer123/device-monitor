package kafkaInterface

import (
	"fmt"
	"sync"
)

func KafkaInterface(waitGroup *sync.WaitGroup) {
	fmt.Println("KafkaInterface started")
	defer waitGroup.Done()
}
