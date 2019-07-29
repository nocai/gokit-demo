package infra

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"nocai/gokit-demo/infra/returncodes"
)

func ServerOptions(l log.Logger) []httptransport.ServerOption {
	return []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(l)),
		httptransport.ServerErrorEncoder(ErrorEncoder),
	}
}

// EncodeResponse 成功时的统一数据处理
func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	switch response.(type) {
	case returncodes.ReturnCoder:
		return httptransport.EncodeJSONResponse(ctx, w, response)
	default:
		return httptransport.EncodeJSONResponse(ctx, w, returncodes.Succ(response))
	}
}

// ErrorEncoder 失败时的统一数据处理
func ErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	switch err.(type) {
	case returncodes.ErrorCoder:
		httptransport.DefaultErrorEncoder(ctx, err, w)
	default:
		httptransport.DefaultErrorEncoder(ctx, returncodes.Fail(err), w)
	}
}
