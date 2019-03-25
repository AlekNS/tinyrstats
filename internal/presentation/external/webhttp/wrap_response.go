package webhttp

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

// encodeResponseWrapData wraps response to sub data object.
func encodeResponseWrapData(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return httptransport.EncodeJSONResponse(ctx, w, map[string]interface{}{
		"data": response,
	})
}
