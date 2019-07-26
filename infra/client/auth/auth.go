package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"net/url"
	"strings"
)

type Service interface {
	Verify()
	Ping()
}

type Endpoints struct {
	PingEndpoint   endpoint.Endpoint
	VerifyEndpoint endpoint.Endpoint
}

func (ep Endpoints) Verify() {

}

func (ep Endpoints) Ping() {
	resp, e := ep.PingEndpoint(context.Background(), nil)
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(resp)
}

func MakePingEndpoint(ser Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		ser.Ping()
		return nil, nil
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
	var response string
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}
