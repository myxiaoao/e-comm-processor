package nats

import (
	"github.com/nats-io/nats-server/v2/server"
)

func startTestServer() (*server.Server, error) {
	opts := &server.Options{
		Host: "127.0.0.1",
		Port: -1, // 随机端口
	}
	ns, err := server.NewServer(opts)
	if err != nil {
		return nil, err
	}
	go ns.Start()
	if !ns.ReadyForConnections(5e9) { // 5秒超时
		return nil, err
	}
	return ns, nil
}
