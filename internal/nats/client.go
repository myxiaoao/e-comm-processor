package nats

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nats.go"

	"ecommerce-processor/internal/domain"
)

type Client struct {
	conn    *nats.Conn
	timeout time.Duration
}

func NewClient(url string, timeout time.Duration) (*Client, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, timeout: timeout}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Conn() *nats.Conn {
	return c.conn
}

func (c *Client) Request(subject string, order domain.Order) (domain.ServiceResponse, error) {
	reqData, err := json.Marshal(order)
	if err != nil {
		return domain.ServiceResponse{}, err
	}

	msg, err := c.conn.Request(subject, reqData, c.timeout)
	if err != nil {
		return domain.ServiceResponse{}, err
	}

	var resp domain.ServiceResponse
	if err := json.Unmarshal(msg.Data, &resp); err != nil {
		return domain.ServiceResponse{}, err
	}

	if !resp.Success {
		return resp, errors.New(resp.Message)
	}

	return resp, nil
}
