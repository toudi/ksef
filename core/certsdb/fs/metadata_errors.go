package fs

import "errors"

var (
	ErrDecodingCertificates          = errors.New("unable to decode certificates array from file")
	ErrUnknownCertificateUID         = errors.New("unknown certificate UID")
	ErrUIDGenerationAttemptsExceeded = errors.New("UID generation attempts exceeded")
	ErrNIPNotOnTheList               = errors.New("NIP is not on the list of this certificate's NIP's")
)
