package certsdb

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"ksef/internal/runtime"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

var errUnsupportedPrivateKeyType = errors.New("unsupported private key type (not EC PRIVATE KEY)")

type CertificateHash struct {
	Environment runtime.Gateway `yaml:"environment"`
	Usage       []Usage         `yaml:"usage"`
	ValidFrom   time.Time       `yaml:"valid-from,omitempty"`
	ValidTo     time.Time       `yaml:"valid-to,omitempty"`
}

func (ch CertificateHash) UsageAsString() string {
	return strings.Join(lo.Map(ch.Usage, func(u Usage, _ int) string { return string(u) }), ", ")
}

func (ch CertificateHash) Hash() string {
	slices.Sort(ch.Usage)
	return fmt.Sprintf("%s:%s:%d:%d", string(ch.Environment), ch.UsageAsString(), ch.ValidFrom.Unix(), ch.ValidTo.Unix())
}

type Certificate struct {
	CertificateHash `yaml:",inline"`
	available       bool
	removed         bool
	UID             string `yaml:"uid"`
	SelfSigned      bool   `yaml:"self-signed,omitempty"`
	NIP             string `yaml:"nip,omitempty"`
	// only applicable to ksef-issued certs
	ReferenceNumber string `yaml:"ref-no,omitempty"`
	SerialNumber    string `yaml:"serial-number,omitempty"`
	CSRData         string `yaml:"csr-data,omitempty"`
	ProfileName     string `yaml:"profile,omitempty"`
}

func (c Certificate) Filename() string {
	return filepath.Join(path.Dir(certificatesDBFile), string(c.Environment)+"-"+c.UID+".pem")
}

func (c Certificate) PrivateKeyFilename() string {
	return filepath.Join(filepath.Dir(certificatesDBFile), string(c.Environment)+"-"+c.UID+"-pkey.pem")
}

func (c Certificate) Expired() bool {
	if c.ValidFrom.IsZero() {
		return false
	}
	now := time.Now()
	return c.ValidTo.Before(now)
}

func (c *Certificate) SavePKey(privateKey *ecdsa.PrivateKey) error {
	privateKeyFile, err := os.Create(c.PrivateKeyFilename())
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	return pem.Encode(privateKeyFile, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privateKeyBytes})
}

func (c *Certificate) SaveCert(derBytes []byte) error {
	certFile, err := os.Create(c.Filename())
	if err != nil {
		return err
	}
	defer certFile.Close()

	return pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})
}

func (c *Certificate) SignContent(content []byte) ([]byte, error) {
	// we can only do it if we have a private key
	pkeyFilename := c.PrivateKeyFilename()
	privKeyBytes, err := os.ReadFile(pkeyFilename)
	if err != nil {
		return nil, err
	}
	privKeyBlock, _ := pem.Decode(privKeyBytes)
	ecPrivKeyAny, err := x509.ParsePKCS8PrivateKey(privKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	ecPrivKey, ok := ecPrivKeyAny.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errUnsupportedPrivateKeyType
	}

	contentHash := sha256.Sum256(content)
	r, s, err := ecdsa.Sign(rand.Reader, ecPrivKey, contentHash[:])
	if err != nil {
		return nil, err
	}
	digest := r.Bytes()
	digest = append(digest, s.Bytes()...)
	return digest, nil
}
