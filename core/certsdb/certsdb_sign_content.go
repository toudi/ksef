package certsdb

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"strings"
)

// SignContent signs the given content with the private key addressed by the
// given key ID (a certificate's PrivateKeyId). Therefore it is up to the
// caller to first retrieve certificate and then retrieve the certificate's
// PrivateKeyId. This may seem weird, however due to the usage of XaDES, the
// signature also contains certificate metadata therefore you have to retrieve
// the certificate meta regardless.
func (c *CertsDB) SignContent(ctx context.Context, privateKeyId string, content []byte) ([]byte, error) {
	if privateKeyId == "" {
		return nil, ErrNoPrivateKey
	}

	privKeyBytes, err := c.dataStorage.GetPrivateKeyBytes(ctx, privateKeyId)
	if err != nil {
		return nil, err
	}
	privKeyBlock, _ := pem.Decode(privKeyBytes)
	var ecPrivKey *ecdsa.PrivateKey
	ecPrivKeyAny, err := x509.ParsePKCS8PrivateKey(privKeyBlock.Bytes)
	if err != nil {
		if strings.Contains(err.Error(), "use ParseECPrivateKey") {
			ecPrivKey, err = x509.ParseECPrivateKey(privKeyBlock.Bytes)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		var ok bool
		ecPrivKey, ok = ecPrivKeyAny.(*ecdsa.PrivateKey)
		if !ok {
			return nil, ErrUnsupportedPrivateKeyType
		}
	}

	contentHash := sha256.Sum256(content)
	r, s, err := ecdsa.Sign(rand.Reader, ecPrivKey, contentHash[:])
	if err != nil {
		return nil, err
	}

	// a fix for situations where either r or s would have leading
	// 0 bytes which would yield a shorter final buffer.
	// with FillBytes we're always filling up to the whole content,
	// padding with zero bytes if required.
	// also, for reference:
	// https://github.com/golang-jwt/jwt/blob/main/ecdsa.go#L127
	curveByteSize := (ecPrivKey.Curve.Params().BitSize + 7) / 8
	digest := make([]byte, 2*curveByteSize)
	r.FillBytes(digest[:curveByteSize])
	s.FillBytes(digest[curveByteSize:])

	// digest := make([]byte, ecPrivKey.Params().BitSize)
	// digest := r.Bytes()
	// digest = append(digest, s.Bytes()...)
	return digest, nil
}
