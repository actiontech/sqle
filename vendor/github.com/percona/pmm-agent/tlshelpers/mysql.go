// pmm-agent
// Copyright 2019 Percona LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package tlshelpers contains helpers for databases tls connections.
package tlshelpers

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// RegisterMySQLCerts is used for register TLS config before sql.Open is called.
func RegisterMySQLCerts(files map[string]string) error {
	if files == nil {
		return nil
	}

	ca := x509.NewCertPool()
	cert, err := tls.X509KeyPair([]byte(files["tlsCert"]), []byte(files["tlsKey"]))
	if err != nil {
		return errors.Wrap(err, "register MySQL client cert failed")
	}

	if ok := ca.AppendCertsFromPEM([]byte(files["tlsCa"])); ok {
		err = mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:      ca,
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			return errors.Wrap(err, "register MySQL CA cert failed")
		}
	}

	return nil
}

// DeregisterMySQLCerts is used for deregister TLS config.
func DeregisterMySQLCerts() {
	mysql.DeregisterTLSConfig("custom")
}
