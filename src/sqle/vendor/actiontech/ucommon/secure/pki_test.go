package secure

import (
	"crypto/x509"
	"encoding/pem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateSecretKey(t *testing.T) {
	caCertPem, caKeyPem, err := GenerateCaCertificate()
	assert.NoError(t, err)
	assert.NotEqual(t, "", caCertPem, "caCert shouldn't be empty")
	assert.NotEqual(t, "", caKeyPem, "caKey shouldn't be empty")

	block, _ := pem.Decode([]byte(caCertPem))
	assert.True(t, block.Type == "CERTIFICATE")
	caCert, err := x509.ParseCertificate(block.Bytes)

	block, _ = pem.Decode([]byte(caKeyPem))
	assert.True(t, block.Type == "RSA PRIVATE KEY")
	caKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	certPem, keyPem, err := GenerateHostCertificate([]string{"127.0.0.1"}, caCert, caKey)
	assert.NoError(t, err)
	assert.NotEqual(t, "", certPem, "cert shouldn't be empty")
	assert.NotEqual(t, "", keyPem, "key shouldn't be empty")

	block, _ = pem.Decode([]byte(certPem))
	assert.True(t, block.Type == "CERTIFICATE")
	cert, err := x509.ParseCertificate(block.Bytes)
	assert.NotNil(t, cert, "cert shouldn't be nil")

}
