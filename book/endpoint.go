package book

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeBookEndpoint(bs Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return bs.Books()
	}
}

func makeGetByIdEndpoint(ser Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (response interface{}, err error) {
		book, err := ser.GetById(request.(int64))
		if err != nil {
			return nil, err
		}
		return book, nil
	}
}
