package certsdb

import (
	"context"
	"iter"
)

// CertificateMetadataStorage handles certificate metadata only: the CRU (no D)
// operations, NIP's, internal ID's and updates. It never reads or writes raw
// bytes (certificate PEM files, private keys) - those live in the
// CertificateDataStorage. Yes, I do realize that this may seem like overenginnering
// (since everything is still inside the `certificates` directory)
// but I do have a reason behind this and that is that I want to achieve two things
// 1/ be backwards compatible with the filesystem-based implementation
// 2/ to make the integration easier. so essentially everything that *can* be
// abstracted from the filesystem implementation can now use this interface to
// interact with CertsDB
type CertificateMetadataStorage interface {
	Certificates(ctx context.Context) ([]Certificate, error)
	// GetCertificate performs a single lookup by certificate UID.
	GetCertificate(ctx context.Context, certId string) (*Certificate, error)
	// UpdateCertificate applies the metadata update. DERBytes is not written
	// by this implementation - the raw bytes belong to the data storage.
	UpdateCertificate(ctx context.Context, certId string, updateData CertificateUpdateData) error
	// StoreCSRRequest stores a CSR request and returns the assigned certificate UID.
	StoreCSRRequest(ctx context.Context, csrData CSRRequestData) (string, error)
	UpsertSelfSignedCertificate(ctx context.Context, selfSignedCertData UserCertificateData) error
	// the implementation should **NOT** store duplicate if a certificate can be found with the given hash
	StoreUserCertificate(ctx context.Context, certData UserCertificateData) (bool, error)
	StoreMinistryCertificate(ctx context.Context, certData MinistryCertificateData) (bool, error)

	// internal id's related functions
	AddInternalIDsToCertificate(ctx context.Context, certId string, internalIds iter.Seq[string]) error
	RemoveInternalIDsFromCertificate(ctx context.Context, certId string, internalIds iter.Seq[string]) error

	// NIP management functions
	AddNIPToCertificate(ctx context.Context, certId string, NIP string) error
	RemoveNIPFromCertificate(ctx context.Context, certId string, NIP string) error

	// implementation may use this method whenever it wants to sync the internal state to persistance layer
	Close(ctx context.Context) error
}

// CertificateDataStorage handles raw bytes only: certificate PEM files and
// private keys, addressed by certificate / key ID.
type CertificateDataStorage interface {
	StoreCertificateBytes(ctx context.Context, certId string, derBytes []byte) error
	GetCertificateBytes(ctx context.Context, certId string) ([]byte, error)

	// StorePrivateKeyBytes creates-or-overwrites the private key bytes under the given key ID.
	StorePrivateKeyBytes(ctx context.Context, keyId string, pemBytes []byte) error
	GetPrivateKeyBytes(ctx context.Context, keyId string) ([]byte, error)

	Close(ctx context.Context) error
}
