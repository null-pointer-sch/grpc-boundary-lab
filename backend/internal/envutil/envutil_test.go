package envutil_test

import (
	"os"
	"testing"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/envutil"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Clear all env vars for this test
	os.Clearenv()

	cfg := envutil.LoadConfig()
	assert.Equal(t, "50051", cfg.BackendPort)
	assert.Equal(t, "50151", cfg.BackendPortTLS)
	assert.Equal(t, "/certs", cfg.CertDir)
	assert.Equal(t, "127.0.0.1", cfg.BackendHost)

	// Set some env vars
	os.Setenv("BACKEND_PORT", "9999")
	os.Setenv("CERT_DIR", "/tmp/certs")
	defer os.Clearenv()

	cfg = envutil.LoadConfig()
	assert.Equal(t, "9999", cfg.BackendPort)
	assert.Equal(t, "/tmp/certs", cfg.CertDir)
}

func TestGetOrDefault(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, "fallback", envutil.GetOrDefault("TEST_VAR", "fallback"))
	
	os.Setenv("TEST_VAR", "actual")
	defer os.Clearenv()
	assert.Equal(t, "actual", envutil.GetOrDefault("TEST_VAR", "fallback"))
}

func TestGetInt(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, 42, envutil.GetInt("TEST_INT", 42))

	os.Setenv("TEST_INT", "100")
	defer os.Clearenv()
	assert.Equal(t, 100, envutil.GetInt("TEST_INT", 42))
}

func TestGetInt64(t *testing.T) {
	os.Clearenv()
	assert.Equal(t, int64(42), envutil.GetInt64("TEST_INT64", 42))

	os.Setenv("TEST_INT64", "100")
	defer os.Clearenv()
	assert.Equal(t, int64(100), envutil.GetInt64("TEST_INT64", 42))
}
