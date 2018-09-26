package secure

import (
	"actiontech/ucommon/os"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"sync"
	"actiontech/ucommon/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"time"
)

// consul key definition
const (
	CLUSTER_SECURITY_KEY   = "universe/security"
	CA_CERTIFICATE_KEY     = "universe/secret-key/certificate"
	CA_PRIVATE_KEY         = "universe/secret-key/privatekey"
	SERVER_CERTIFICATE_KEY = "universe/servers/%v/secret-key/certificate"
	SERVER_PRIVATE_KEY     = "universe/servers/%v/secret-key/privatekey"
	SERVER_ADDR_KEY        = "universe/servers/%v/addr"
)

// pem file store path
const (
	CA_CERTIFICATE_PATH     = "/etc/pki/actiontech-universe/ca_cert.pem"
	CA_PRIVATE_PATH         = "/etc/pki/actiontech-universe/ca_key.pem"
	SERVER_CERTIFICATE_PATH = "/etc/pki/actiontech-universe/cert.pem"
	SERVER_PRIVATE_PATH     = "/etc/pki/actiontech-universe/key.pem"
)

var _securityEnabled bool = false
var _securityEnabledMutex sync.Mutex

func SetSecurityEnabled(securityEnabled bool) {
	_securityEnabledMutex.Lock()
	_securityEnabled = securityEnabled
	_securityEnabledMutex.Unlock()

	stage := log.NewStage().Enter("security")
	if securityEnabled {
		log.Key(stage, "security enabled")
	} else {
		log.Key(stage, "security disabled")
	}
}

func IsSecurityEnabled() bool {
	_securityEnabledMutex.Lock()
	defer _securityEnabledMutex.Unlock()
	return _securityEnabled
}

//shortcut
func NewGrpcServer(opt ...grpc.ServerOption) (*grpc.Server, error) {
	keepaliveConf := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime: 1 * time.Second,
		PermitWithoutStream: true,
	})
	if IsSecurityEnabled() {
		creds, err := GetServerTLSCredentials()
		if nil != err {
			return nil, err
		}
		return grpc.NewServer(append(opt, grpc.Creds(creds), keepaliveConf)...), nil
	} else {
		return grpc.NewServer(append(opt, keepaliveConf)...), nil
	}
}

var _transportCredentials credentials.TransportCredentials
var _transportCredentialsMu sync.Mutex

func GetClientTLSCredentials() (credentials.TransportCredentials, error) {
	_transportCredentialsMu.Lock()

	tc := _transportCredentials
	if nil != tc {
		_transportCredentialsMu.Unlock()
		return tc, nil
	}

	rawCACert, err := ioutil.ReadFile(CA_CERTIFICATE_PATH)
	if nil != err {
		_transportCredentialsMu.Unlock()
		return nil, fmt.Errorf("get ca certificate failed: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(rawCACert)
	new := credentials.NewTLS(&tls.Config{
		RootCAs: caCertPool,
	})
	_transportCredentials = new
	tc = _transportCredentials
	_transportCredentialsMu.Unlock()
	return tc, nil
}

func GetServerTLSCredentials() (credentials.TransportCredentials, error) {
	if !os.IsFileExist(SERVER_CERTIFICATE_PATH) {
		return nil, fmt.Errorf("server certificate not found")
	}
	if !os.IsFileExist(SERVER_PRIVATE_PATH) {
		return nil, fmt.Errorf("server privatekey not found")
	}
	cert, err := tls.LoadX509KeyPair(SERVER_CERTIFICATE_PATH, SERVER_PRIVATE_PATH)
	if err != nil {
		return nil, fmt.Errorf("get ca certificate failed: %v", err)
	}
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
	}), nil
}

func GetCaCertificate() (certificate string, err error) {
	certificate, err = os.GetFileContent(CA_CERTIFICATE_PATH)
	if err != nil {
		return "", fmt.Errorf("fail to get ca certificate: %v", err)
	}
	return certificate, nil
}

func GetCaPrivateKey() (privateKey string, err error) {
	privateKey, err = os.GetFileContent(CA_PRIVATE_PATH)
	if err != nil {
		return "", fmt.Errorf("fail to get ca private key: %v", err)
	}
	return privateKey, nil
}

func GetServerCertificate() (certificate string, err error) {
	certificate, err = os.GetFileContent(SERVER_CERTIFICATE_PATH)
	if err != nil {
		return "", fmt.Errorf("fail to get server certificate: %v", err)
	}
	return certificate, nil
}

func GetServerPrivateKey() (privateKey string, err error) {
	privateKey, err = os.GetFileContent(CA_PRIVATE_PATH)
	if err != nil {
		return "", fmt.Errorf("fail to get server private key: %v", err)
	}
	return privateKey, nil
}
