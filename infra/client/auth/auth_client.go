package auth

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"io"
	"time"
)

func New(l log.Logger, consulClient consul.Client) Service {
	var (
		consulService = "ivargo-auth"
		consulTags    = []string{}
		passingOnly   = true
		retryMax      = 3
		retryTimeout  = 500 * time.Millisecond
	)

	var (
		instancer = consul.NewInstancer(consulClient, l, consulService, consulTags, passingOnly)
		endpoints Endpoints
	)

	{
		factory := factoryFor(MakePingEndpoint)
		endpointer := sd.NewEndpointer(instancer, factory, l)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(retryMax, retryTimeout, balancer)
		endpoints.PingEndpoint = retry
	}

	return endpoints
}

func factoryFor(makeEndpoint func(Service) endpoint.Endpoint) sd.Factory {
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {

		service, err := MakeClientEndpoints(instance)
		if err != nil {
			return nil, nil, err
		}
		return makeEndpoint(service), nil, nil
	}
}
