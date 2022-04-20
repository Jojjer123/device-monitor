package streamManager

import (
	"fmt"
	"sync"

	"github.com/onosproject/monitor-service/pkg/types"
	// reqBuilder "github.com/onosproject/monitor-service/pkg/requestBuilder"
)

// var streamStore []types.Stream

func StreamManager(waitGroup *sync.WaitGroup, streamMgrChannel chan types.StreamMgrChannelMessage) { //, adminChannel chan types.AdminChannelMessage) {
	defer waitGroup.Done()

	// TODO: Remove streamWaitGroup and add better way of keeping module "alive".

	var streamWaitGroup sync.WaitGroup

	var streamMgrMessage types.StreamMgrChannelMessage
	streamMgrMessage.ManageCmd = streamMgrCmd
	streamMgrChannel <- streamMgrMessage

	streamWaitGroup.Wait()
}

func streamMgrCmd(stream types.Stream, cmd string) string {
	fmt.Println(stream.Target, cmd)

	return ""
}
