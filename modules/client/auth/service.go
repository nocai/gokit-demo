package auth

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"nocai/gokit-demo/infra/returncodes"
)

type Service interface {
	Verify()
	Ping() (returncodes.ReturnCoder, error)
}

type Endpoints struct {
	PingEndpoint   endpoint.Endpoint
	VerifyEndpoint endpoint.Endpoint
}

func (ep Endpoints) Verify() {

}

func (ep Endpoints) Ping() (returncodes.ReturnCoder, error) {
	resp, err := ep.PingEndpoint(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return resp.(returncodes.ReturnCoder), nil
}
