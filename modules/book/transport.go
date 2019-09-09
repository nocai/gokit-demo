package book

import (
	"context"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/nocai/gokit-demo/infra"
	"github.com/nocai/gokit-demo/infra/returncodes"
	"strconv"
)

func MakeHandler(bs Service, l log.Logger) http.Handler {
	opts := infra.ServerOptions(l)

	booksHandler := httptransport.NewServer(
		makeBookEndpoint(bs),
		httptransport.NopRequestDecoder,
		infra.EncodeResponse,
		opts...,
	)

	getByIdHandler := httptransport.NewServer(
		makeGetByIdEndpoint(bs),
		decodeGetByIdRequest,
		infra.EncodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/books/", booksHandler)
	r.Handle("/books/{id}", getByIdHandler)
	return r
}

func decodeGetByIdRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, returncodes.ErrBadRequest
	}

	id64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	return id64, nil
}
