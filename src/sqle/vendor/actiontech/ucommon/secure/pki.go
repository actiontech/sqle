package secure

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

var RSA_KEY_LEN = 1024

// generate a ca certificate which is self-signed
func GenerateCaCertificate() (certificate string, privateKey string, error error) {
	key, err := rsa.GenerateKey(rand.Reader, RSA_KEY_LEN)
	if err != nil {
		return "", "", err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(100 * 365 * 24 * time.Hour))

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Action Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return "", "", err
	}
	buf := bytes.Buffer{}
	if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); nil != err {
		return "", "", err
	}
	certificate = buf.String()
	buf.Reset()
	if err := pem.Encode(&buf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); nil != err {
		return "", "", err
	}
	privateKey = buf.String()
	return certificate, privateKey, nil
}

// generate a host certificate which is ca-signed, host either a hostname or ip address that is binding to the certificate
func GenerateHostCertificate(hosts []string, caCertificate *x509.Certificate, caPrivateKey *rsa.PrivateKey) (certificate string, privateKey string, error error) {
	key, err := rsa.GenerateKey(rand.Reader, RSA_KEY_LEN)
	if err != nil {
		return "", "", err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(10 * 365 * 24 * time.Hour))

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", "", err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Action Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, host := range hosts {
		if ip := net.ParseIP(host); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, host)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCertificate, &key.PublicKey, caPrivateKey)
	if err != nil {
		return "", "", err
	}
	buf := bytes.Buffer{}
	if err := pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); nil != err {
		return "", "", err
	}
	certificate = buf.String()
	buf.Reset()
	if err := pem.Encode(&buf, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); nil != err {
		return "", "", err
	}
	privateKey = buf.String()
	return certificate, privateKey, nil
}
