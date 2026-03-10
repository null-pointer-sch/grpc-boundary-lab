package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/httputil"
	pb "github.com/null-pointer-sch/grpc-boundary-lab/internal/proto"
	"github.com/null-pointer-sch/grpc-boundary-lab/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlePing_Get(t *testing.T) {
	h := &service.RESTPingHandler{}

	req := httptest.NewRequest(http.MethodGet, "/ping?message=hello", nil)
	w := httptest.NewRecorder()

	h.HandlePing(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var payload pb.PingResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	assert.Equal(t, "pong: hello", payload.Message)
}

func TestHandlePing_Post(t *testing.T) {
	h := &service.RESTPingHandler{}

	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	w := httptest.NewRecorder()

	h.HandlePing(w, req)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)

	var payload httputil.ErrorResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	assert.Equal(t, "Method not allowed", payload.Message)
}
