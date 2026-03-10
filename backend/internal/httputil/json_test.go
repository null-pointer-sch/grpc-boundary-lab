package httputil_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/httputil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}

	httputil.WriteJSON(w, data)

	res := w.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	var payload map[string]string
	require.NoError(t, json.NewDecoder(res.Body).Decode(&payload))
	assert.Equal(t, "bar", payload["foo"])
}

// badWriter is a mock ResponseWriter that fails on Write.
type badWriter struct {
	http.ResponseWriter
}

func (w *badWriter) Write(b []byte) (int, error) {
	return 0, http.ErrBodyNotAllowed
}

type unencodable struct {
	Ch chan int // channels cannot be JSON encoded
}

func TestWriteJSON_Error(t *testing.T) {
	w := httptest.NewRecorder()

	// attempt to encode something un-encodable to trigger the true error path
	httputil.WriteJSON(w, unencodable{Ch: make(chan int)})

	res := w.Result()
	defer res.Body.Close()

	// should write the error out natively
	// error output will just be text with a 500 status code
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
