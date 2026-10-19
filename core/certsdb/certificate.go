package certsdb

import (
	"ksef/core/ksef"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

// userCertificateUsages are the usages a user may import their own certificate
// for. The other usages (KsefTokenEncryption, SymmetricKeyEncryption) are issued
// by the ministry of finance and enter the database via StoreMinistryCertificate only.
var userCertificateUsages = []CertificateUsage{
	CertificateUsageAuthentication,
	CertificateUsageOffline,
}

type CertificateHash struct {
	Gateway   ksef.KSeFGateway
	Usage     []CertificateUsage
	ValidFrom time.Time
	ValidTo   time.Time
}

func (c CertificateHash) UsageAsString() string {
	return strings.Join(lo.Map(c.Usage, func(u CertificateUsage, _ int) string { return string(u) }), ", ")
}

type Certificate struct {
	CertificateHash
	UID          string
	SerialNumber string
	NIP          []string
	ProfileName  string
	SelfSigned   bool
	CN           *string
	// KSeF certificate is combined with privateKey
	PrivateKeyId    string
	InternalIDs     []string
	CSRData         string
	ReferenceNumber string
}

// Validate checks the properties common to certificates of every origin
// (user-imported or ministry-issued).
func (c Certificate) Validate() error {
	if len(c.Usage) == 0 {
		return ErrInvalidCertificateUsage
	}
	return nil
}

// ValidateUserCertificate extends Validate with the rules for certificates a
// user imports themselves (in contrast to ministry-issued ones).
func (c Certificate) ValidateUserCertificate() error {
	if err := c.Validate(); err != nil {
		return err
	}
	for _, usage := range c.Usage {
		if !slices.Contains(userCertificateUsages, usage) {
			return ErrInvalidCertificateUsage
		}
	}
	return nil
}

// CertificateUpdateData holds potential updates to the certificate
// - hence everything is optional here. the metadata implementation
// interprets this structure and updates certificate's properties
// accordingly.
type CertificateUpdateData struct {
	SerialNumber    *string
	ReferenceNumber *string
	ValidFrom       *time.Time
	ValidTo         *time.Time
	DERBytes        []byte
	ProfileName     *string
}

type MinistryCertificateData struct {
	Certificate Certificate
	DERBytes    []byte
}

type UserCertificateData struct {
	PrivateKeyPEMBytes []byte
	Certificate        Certificate
	DERBytes           []byte
}

// a CSR stands for Certificate Signing Request. Basically a customer
// can generate it, then send it to KSeF - along with private key.
// KSeF will process it (it's an asynchronous process) and generate a
// certificate file.
type CSRRequestData struct {
	Certificate        Certificate
	PrivateKeyPEMBytes []byte
}

func (c Certificate) IsEncryptable() bool {
	// only these two types are actually encryptable (i.e. they contain private key)
	return slices.ContainsFunc(c.Usage, func(u CertificateUsage) bool {
		return u == CertificateUsageAuthentication || u == CertificateUsageOffline
	})
}
