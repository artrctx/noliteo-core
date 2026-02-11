package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"testing"

	"github.com/artrctx/noliteo-core/internal/config"
)

func randPrivateKey(bits int) (string, error) {
	// 1. Generate the RSA private key.
	// Use crypto/rand.Reader for cryptographically secure randomness.
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// 2. Marshal the private key to PKCS#8, ASN.1 DER form.
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)

	if err != nil {
		return "", fmt.Errorf("failed to generate RSA PKCS#8 key: %w", err)
	}

	// 3. Encode the DER bytes into a PEM block.
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	// 4. Encode the PEM block to a byte slice and convert to a string.
	privateKeyPem := pem.EncodeToMemory(pemBlock)
	return string(privateKeyPem), nil
}

func TestMain(m *testing.M) {
	// generate RSA private key
	test_pk, err := randPrivateKey(2048)
	if err != nil {
		log.Fatalf("failed generating test 32 bytes. err: %v", err)
	}

	jwtMgr = &JwtManager{config.NewJwtConfig(test_pk)}

	m.Run()
}

func TestEncodeAndDecodeJwt(t *testing.T) {
	cases := [5]Token{
		{"id1", rand.Text()},
		{"id2", rand.Text()},
		{"id3", rand.Text()},
		{"id4", rand.Text()},
		{"id5", rand.Text()},
	}

	for _, tkn := range cases {
		jwt, err := GenerateToken(tkn)
		if err != nil {
			t.Errorf("failed generating token with token ident: %v, tid: %v; err: %v", tkn.TID, tkn.Ident, err)
		}

		jTkn, err := ValidateToken(jwt)
		if err != nil {
			t.Errorf("failed to verify token; err: %v", err)
		}

		if jTkn.TID != tkn.TID {
			t.Errorf("parsed jwt TID different from encoded value; expected %v; got: %v;", tkn.TID, jTkn.TID)
		}

		if jTkn.Ident != tkn.Ident {
			t.Errorf("parsed jwt Ident different from encoded value; expected %v; got: %v;", tkn.Ident, jTkn.Ident)
		}
	}

}
