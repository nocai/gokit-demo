package infra

import (
	"context"
	"encoding/json"
	"net/http"
	"nocai/gokit-demo/infra/returncodes"
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	switch response.(type) {
	case returncodes.ReturnCoder:
		return json.NewEncoder(w).Encode(response)
	case error:
		return json.NewEncoder(w).Encode(returncodes.Fail(response))
	default:
		return json.NewEncoder(w).Encode(returncodes.Succ(response))
	}
}

