package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/pkcs12"
)

func main() {

	fqdn := os.Args[1]
	secret := os.Args[2]

	certData, err := os.ReadFile(fqdn + ".p12")
	if err != nil {
		log.Fatalln("failed", err)
	}

	certificate, rsaPrivateKey, err := decodePkcs12(certData, secret)
	if err != nil {
		log.Fatalln("failed", err)
	}

	key, _ := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)

	cert := certificate.Raw

	keyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: key,
	}

	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}

	cf, err := os.Create(fqdn + "-signed.crt")
	if err != nil {
		panic(err)
	}
	defer cf.Close()
	if err := pem.Encode(cf, certBlock); err != nil {
		log.Fatal(err)
	}

	kf, err := os.Create(fqdn + "-key.pem")
	if err != nil {
		panic(err)
	}
	defer kf.Close()
	if err := pem.Encode(kf, keyBlock); err != nil {
		log.Fatal(err)
	}
}

func decodePkcs12(pkcs []byte, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, certificate, err := pkcs12.Decode(pkcs, password)
	if err != nil {
		return nil, nil, err
	}

	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey {
		return nil, nil, fmt.Errorf("PKCS#12 certificate must contain an RSA private key")
	}

	return certificate, rsaPrivateKey, nil
}
