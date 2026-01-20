package nats

import (
	"errors"

	"github.com/nats-io/nats-server/v2/server"
)

var ErrServerNotReady = errors.New("NATS server failed to start within timeout")

// StartTestServer creates an embedded NATS server for testing.
// Returns a server listening on a random port.
func StartTestServer() (*server.Server, error) {
	opts := &server.Options{
		Host: "127.0.0.1",
		Port: -1, // random port
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	go ns.Start()
	if !ns.ReadyForConnections(5e9) { // 5 second timeout
		return nil, ErrServerNotReady
	}
	return ns, nil
}
