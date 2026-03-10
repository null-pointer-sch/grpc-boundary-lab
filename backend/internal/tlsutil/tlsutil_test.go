package tlsutil_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/null-pointer-sch/grpc-boundary-lab/internal/tlsutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func generateTestCertAndKey(t *testing.T) (string, string) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		IsCA:      true,
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	require.NoError(t, err)

	certFile, err := os.CreateTemp("", "cert-*.pem")
	require.NoError(t, err)
	defer certFile.Close()
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	keyFile, err := os.CreateTemp("", "key-*.pem")
	require.NoError(t, err)
	defer keyFile.Close()
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	require.NoError(t, err)
	pem.Encode(keyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	return certFile.Name(), keyFile.Name()
}

func TestLoadServerConfig(t *testing.T) {
	certFile, keyFile := generateTestCertAndKey(t)
	defer os.Remove(certFile)
	defer os.Remove(keyFile)

	cfg, err := tlsutil.LoadServerConfig(certFile, keyFile)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.Len(t, cfg.Certificates, 1)

	// Test error case
	_, err = tlsutil.LoadServerConfig("nonexistent-cert.pem", "nonexistent-key.pem")
	assert.Error(t, err)
}

func TestLoadClientConfig(t *testing.T) {
	certFile, keyFile := generateTestCertAndKey(t)
	defer os.Remove(certFile)
	defer os.Remove(keyFile)

	cfg, err := tlsutil.LoadClientConfig(certFile)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	assert.NotNil(t, cfg.RootCAs)

	// Test error case
	_, err = tlsutil.LoadClientConfig("nonexistent-ca.pem")
	assert.Error(t, err)
}

func TestLoadClientConfig_InvalidCert(t *testing.T) {
	badCertFile, err := os.CreateTemp("", "badcert-*.pem")
	require.NoError(t, err)
	defer os.Remove(badCertFile.Name())
	
	badCertFile.WriteString("NOT A REAL CERTIFICATE")
	badCertFile.Close()

	_, err = tlsutil.LoadClientConfig(badCertFile.Name())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to append CA certificate to pool")
}
