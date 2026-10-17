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

func (c *CertsDB) SignContent(ctx context.Context, certId string, content []byte) ([]byte, error) {
	// we can only do it if we have a private key
	privateKeyId, err := c.storage.GetPrivateKeyId(ctx, certId)
	if err != nil {
		return nil, err
	}
	privKeyBytes, err := c.storage.GetPrivateKeyBytes(ctx, privateKeyId)
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
