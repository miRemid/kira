package hystrix

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	name := req.Service() + "." + req.Endpoint()
	return hystrix.Do(name,
		func() error {
			return c.Client.Call(ctx, req, rsp, opts...)
		},
		func(err error) error {
			return err
		},
	)
}

func NewClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}
