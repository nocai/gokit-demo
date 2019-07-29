package auth

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/http"
	"net/url"
	"nocai/gokit-demo/infra/returncodes"
	"strings"
	"time"
)

func NewClient(l log.Logger, consulClient consul.Client) Service {
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

func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path. That's fine: we simply need to provide specific encoders for
	// each endpoint.

	return Endpoints{

		PingEndpoint: httptransport.NewClient("GET", tgt, encodePingRequest, decodePingResponse, options...).Endpoint(),
		//VerifyEndpoint: httptransport.NewClient("GET", tgt, encodeGetProfileRequest, decodeGetProfileResponse, options...).Endpoint(),
	}, nil
}

func encodePingRequest(_ context.Context, req *http.Request, _ interface{}) error {
	// r.Methods("POST").Path("/profiles/")
	req.URL.Path = "/ping"
	return nil
}

func decodePingResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	return returncodes.Unmarshal(resp.Body)
}