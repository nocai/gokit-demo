package book

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"nocai/gokit-demo/infra"
	"nocai/gokit-demo/infra/returncodes"
	"strconv"
)

func MakeHandler(bs Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	booksHandler := kithttp.NewServer(
		makeBookEndpoint(bs),
		kithttp.NopRequestDecoder,
		infra.EncodeResponse,
		opts...,
	)

	getByIdHandler := kithttp.NewServer(
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
