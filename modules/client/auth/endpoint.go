package auth

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func MakePingEndpoint(ser Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return ser.Ping()
	}
}
