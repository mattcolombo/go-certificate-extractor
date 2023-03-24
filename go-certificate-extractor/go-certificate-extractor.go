package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/pkcs12"
)

func main() {

	// inputs required for the tool to work
	filepath := os.Args[1]
	fqdn := os.Args[2]
	secret := os.Args[3]

	// reading the certificate from the file provided
	certData, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalln("Reading p12 file for extraction failed", err)
	}

	// decoding the data from the p12 file provided and storing in a certificate and key object
	certificate, rsaPrivateKey, err := decodePkcs12(certData, secret)
	if err != nil {
		log.Fatalln("Extraction of PEM data from p12 failed", err)
	}

	// unmarshal the RSA private key to be readable as an unencrypted PEM key
	key, _ := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)
	// extracting the actual raw certificate information from the certificate object for later conversion to PEM
	cert := certificate.Raw

	//creating a file with the PEM block related to the private key
	keyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: key,
	}
	kf, err := os.Create(fqdn + "-key.pem")
	if err != nil {
		log.Fatalln(err)
	}
	defer kf.Close()
	if err := pem.Encode(kf, keyBlock); err != nil {
		log.Fatalln(err)
	}
	// adding b64 encoding of the key
	createB64EncodedFile(fqdn + "-key.pem")

	// creating a file with the PEM block related to the certificate
	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}
	cf, err := os.Create(fqdn + "-signed.crt")
	if err != nil {
		log.Fatalln(err)
	}
	defer cf.Close()
	if err := pem.Encode(cf, certBlock); err != nil {
		log.Fatalln(err)
	}
	// adding b64 encoding of the above
	createB64EncodedFile(fqdn + "-signed.crt")

}

// this function helps to decode the p12 certificate into a certificate and key object
func decodePkcs12(pkcs []byte, password string) (*x509.Certificate, *rsa.PrivateKey, error) {
	privateKey, certificate, err := pkcs12.Decode(pkcs, password)
	if err != nil {
		return nil, nil, err
	}

	// check that the private key is a valid one. If not return an error
	rsaPrivateKey, isRsaKey := privateKey.(*rsa.PrivateKey)
	if !isRsaKey {
		return nil, nil, fmt.Errorf("PKCS#12 certificate must contain an RSA private key")
	}

	return certificate, rsaPrivateKey, nil
}

func createB64EncodedFile(inputName string) {

	// reading the PEM certificate from the file provided
	pemData, err := os.ReadFile(inputName)
	if err != nil {
		log.Fatalln("Reading file for b64 encoding failed", err)
	}

	b64Data := base64.StdEncoding.EncodeToString(pemData)
	f, err := os.Create(inputName + "-b64enc.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	_, errWrite := f.WriteString("Writing b64 encoded contents of " + inputName + "\n" + b64Data)
	if errWrite != nil {
		log.Fatalln(err)
	}
}
