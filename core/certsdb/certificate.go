package certsdb

import (
	"ksef/core/ksef"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

type Certificate struct {
	UID          string
	Gateway      ksef.KSeFGateway
	SerialNumber string
	NIP          []string
	ProfileName  string
	SelfSigned   bool
	CN           *string
	// KSeF certificate is combined with privateKey
	PrivateKeyId string
	Usage        []CertificateUsage
	ValidFrom    time.Time
	ValidTo      time.Time
}

func (c Certificate) UsageAsString() string {
	return strings.Join(lo.Map(c.Usage, func(u CertificateUsage, _ int) string { return string(u) }), ", ")
}

func (c Certificate) IsEncryptable() bool {
	// only these two types are actually encryptable (i.e. they contain private key)
	return slices.ContainsFunc(c.Usage, func(u CertificateUsage) bool {
		return u == CertificateUsageAuthentication || u == CertificateUsageOffline
	})
}
