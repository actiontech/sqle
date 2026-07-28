package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
)

// Small helper to generate an SM2 PEM private key for login encryption config.
func main() {
	out := "sm2_login_private.pem"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	priv, err := sm2.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	pemBytes, err := x509.WritePrivateKeyToPem(priv, nil)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(out, pemBytes, 0600); err != nil {
		panic(err)
	}
	fmt.Printf("wrote %s\n", out)
}
