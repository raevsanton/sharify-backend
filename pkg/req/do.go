package req

import (
	"fmt"
	"io"
	"net/http"

	"github.com/raevsanton/sharify-backend/pkg/codec"
)

func DoRequest[T any](req *http.Request, expectedStatus int) (T, error) {
	var result T

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != expectedStatus {
		body, _ := io.ReadAll(res.Body)
		return result, fmt.Errorf("received unexpected status: %d, body: %s", res.StatusCode, body)
	}

	result, err = codec.Decode[T](res.Body)
	if err != nil {
		return result, fmt.Errorf("failed to decode response body: %w", err)
	}

	return result, nil
}
