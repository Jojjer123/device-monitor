package types

import "sync"

type AdminChannelMessage struct {
	RegisterFunction func(chan string, *sync.WaitGroup)
	Message          string
}
