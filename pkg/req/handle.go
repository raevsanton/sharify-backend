package req

import (
	"net/http"

	"github.com/raevsanton/sharify-backend/pkg/codec"
	"github.com/raevsanton/sharify-backend/pkg/res"
)

func HandleBody[T any](w *http.ResponseWriter, r *http.Request) (*T, error) {
	defer r.Body.Close()

	body, err := codec.Decode[T](r.Body)
	if err != nil {
		res.Json(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	err = IsValid(body)
	if err != nil {
		res.Json(*w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return &body, nil
}
