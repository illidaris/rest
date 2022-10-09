package signature

import (
	"context"
	"net/http"
	"testing"
)

func TestNewUnSignRequest(t *testing.T) {
	f := SignSetNil.RequestWithContextFunc(nil, nil)
	ctx := context.Background()
	f(ctx, http.MethodGet, "http://localhost:8080", nil)
}
