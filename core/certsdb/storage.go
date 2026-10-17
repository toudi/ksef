package certsdb

import (
	"context"
	"iter"
)

type CertificateStorage interface {
	Certs(ctx context.Context) ([]Certificate, error)
	UpdateCertificate(ctx context.Context, certId string, updateData CertificateUpdateData) error
	StoreCSRRequest(ctx context.Context, cert Certificate) (string, error)
	GetCertificateBytes(ctx context.Context, certId string) ([]byte, error)
	UpsertSelfSignedCertificate(ctx context.Context, selfSignedCertData UserCertificateData) error
	// the implementation should **NOT** store duplicate if a certificate can be found with the given hash
	StoreUserCertificate(ctx context.Context, certData UserCertificateData) (bool, error)
	StoreMinistryCertificate(ctx context.Context, certData MinistryCertificateData) (bool, error)

	// internal id's related functions
	AddInternalIDsToCertificate(ctx context.Context, certId string, internalIds iter.Seq[string]) error
	RemoveInternalIDsFromCertificate(ctx context.Context, certId string, internalIds iter.Seq[string]) error

	// private key related functions
	GetPrivateKeyId(ctx context.Context, certId string) (string, error)
	GetPrivateKeyBytes(ctx context.Context, certId string) ([]byte, error)
	StorePrivateKey(ctx context.Context, pemBytes []byte, gatewayId string) (string, error)

	// NIP management functions
	AddNIPToCertificate(ctx context.Context, certId string, NIP string) error
	RemoveNIPFromCertificate(ctx context.Context, certId string, NIP string) error

	// implementation may use this method whenever it wants to sync the internal state to persistance layer
	Close(ctx context.Context) error
}
