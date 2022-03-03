package adminServer

import (
	"fmt"
	"sync"

	"golang.org/x/net/context"
)

type Server struct {
	// ServerChannel chan string
	RegisterFunction func(chan string, *sync.WaitGroup)
}

func (s *Server) MonitorDevice(ctx context.Context, message *MonitorMessage) (*MonitorResponse, error) {
	fmt.Println("Received message from client: ", message.Action, " ", message.Target)

	if !actionExists(message.Action) {
		// TODO: Set correct error.
		return &MonitorResponse{Response: "Action does not exist."}, nil
	}

	serverChannel := make(chan string)

	var channelWaitGroup sync.WaitGroup
	channelWaitGroup.Add(1)
	go s.RegisterFunction(serverChannel, &channelWaitGroup)

	serverChannel <- message.Action + " " + message.Target

	var response string

	select {
	case response = <-serverChannel:
		{
			fmt.Println("Response from registering function is: ", response)
		}
	}

	channelWaitGroup.Wait()

	return &MonitorResponse{Response: response}, nil
}

func actionExists(action string) bool {
	exists := false
	if action == "Create" || action == "Update" || action == "Delete" {
		exists = true
	}

	return exists
}
