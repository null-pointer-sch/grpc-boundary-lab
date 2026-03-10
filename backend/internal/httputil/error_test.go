package httputil_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("something went wrong")

	httputil.WriteError(w, http.StatusBadRequest, err)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var payload httputil.ErrorResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))

	assert.Equal(t, http.StatusText(http.StatusBadRequest), payload.Error)
	assert.Equal(t, "something went wrong", payload.Message)
}

func TestWriteErrorMessage(t *testing.T) {
	w := httptest.NewRecorder()

	httputil.WriteErrorMessage(w, http.StatusNotFound, "not found here")

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var payload httputil.ErrorResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))

	assert.Equal(t, http.StatusText(http.StatusNotFound), payload.Error)
	assert.Equal(t, "not found here", payload.Message)
}
