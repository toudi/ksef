package certsdb

import "errors"

var (
	ErrUnsupportedPrivateKeyType   = errors.New("unsupported private key type (not EC PRIVATE KEY)")
	ErrCertificateNotFound         = errors.New("could not find matching certificate")
	ErrAmbiguousCertificate        = errors.New("more than one certificate matches, expected exactly one")
	ErrUnsupportedCertificateUsage = errors.New("unsupported certificate usage")
	ErrInvalidCertificateUsage     = errors.New("unexpected certificate usage")
	ErrCertificateAlreadyExists    = errors.New("certificate with this properties is already imported")
)
