package fs

import (
	"fmt"
	"ksef/core/certsdb"
	"ksef/core/ksef"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

type CertificateHash struct {
	EnvironmentId string                     `yaml:"environment"`
	Usage         []certsdb.CertificateUsage `yaml:"usage"`
	ValidFrom     time.Time                  `yaml:"valid-from,omitempty"`
	ValidTo       time.Time                  `yaml:"valid-to,omitempty"`
}

func (ch CertificateHash) UsageAsString() string {
	return strings.Join(lo.Map(ch.Usage, func(u certsdb.CertificateUsage, _ int) string { return string(u) }), ", ")
}

func (ch CertificateHash) Hash() string {
	slices.Sort(ch.Usage)
	return fmt.Sprintf("%s:%s:%d:%d", ch.EnvironmentId, ch.UsageAsString(), ch.ValidFrom.Unix(), ch.ValidTo.Unix())
}

type fsCertificate struct {
	CertificateHash `yaml:",inline"`

	UID        string   `yaml:"uid"`
	SelfSigned bool     `yaml:"self-signed,omitempty"`
	NIPRaw     any      `yaml:"nip,omitempty"`
	NIP        []string `yaml:"-"`
	// only applicable to ksef-issued certs
	ReferenceNumber string  `yaml:"ref-no,omitempty"`
	SerialNumber    string  `yaml:"serial-number,omitempty"`
	CSRData         string  `yaml:"csr-data,omitempty"`
	ProfileName     string  `yaml:"profile,omitempty"`
	CN              *string `yaml:"cn"`
	// https://www.gov.pl/web/finanse/ksef-kiedy-masz-w-firmie-wiecej-pracownikow
	// internal ID has the following format:
	// NIP-12345
	InternalIDs []string `yaml:"internal-ids,omitempty"`

	// TODO: implement marking certificate as removed when it expires
	removed bool
}

func (f fsCertificate) Filename() string {
	return filepath.Join(certificatesDir, f.EnvironmentId+"-"+f.UID+".pem")
}

func (f fsCertificate) AsCertificate() certsdb.Certificate {
	return certsdb.Certificate{
		CertificateHash: certsdb.CertificateHash{
			Gateway:   ksef.KSeFGateway(f.EnvironmentId),
			Usage:     f.Usage,
			ValidFrom: f.ValidFrom,
			ValidTo:   f.ValidTo,
		},
		UID:          f.UID,
		CN:           f.CN,
		SerialNumber: f.SerialNumber,
		SelfSigned:   f.SelfSigned,
		ProfileName:  f.ProfileName,
		NIP:          f.NIP,
		PrivateKeyId: f.UID,

		InternalIDs:     f.InternalIDs,
		ReferenceNumber: f.ReferenceNumber,
		CSRData:         f.CSRData,
	}
}

func fromCertsDBCertificate(input certsdb.Certificate) *fsCertificate {
	return &fsCertificate{
		CertificateHash: CertificateHash{
			EnvironmentId: string(input.Gateway),
			Usage:         input.Usage,
			ValidFrom:     input.ValidFrom,
			ValidTo:       input.ValidTo,
		},
		CN:              input.CN,
		SerialNumber:    input.SerialNumber,
		SelfSigned:      input.SelfSigned,
		ProfileName:     input.ProfileName,
		NIP:             input.NIP,
		InternalIDs:     input.InternalIDs,
		ReferenceNumber: input.ReferenceNumber,
	}
}
