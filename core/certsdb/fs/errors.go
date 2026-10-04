package fs

import "errors"

var (
	ErrDecodingCertificates  = errors.New("unable to decode certificates array from file")
	ErrUnknownCertificateUID = errors.New("unknown certificate UID")
)
