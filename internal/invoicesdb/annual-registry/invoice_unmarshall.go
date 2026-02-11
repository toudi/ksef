package annualregistry

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"ksef/internal/invoice"

	"github.com/fxamacker/cbor/v2"
)

var (
	errDecodingBase64     = errors.New("unable to decode base64 content")
	errConstructingReader = errors.New("unable to construct gzip reader")
	errUncompressing      = errors.New("unable to decompress cbor bytes")
)

func (i *Invoice) Unmarshall() (*invoice.Invoice, error) {
	// the content is serialized in the following way:
	// invoice.Invoice -> CBOR -> gzip -> base64
	// therefore the deserialization needs to happen in the reverse order:
	// base64 -> gunzip -> CBOR -> invoice.Invoice
	gzipBytes, err := base64.StdEncoding.DecodeString(i.Contents)
	if err != nil {
		return nil, errors.Join(errDecodingBase64, err)
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(gzipBytes))
	if err != nil {
		return nil, errors.Join(errConstructingReader, err)
	}
	defer gzipReader.Close()
	cborBytes, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, errors.Join(errUncompressing, err)
	}
	invoice := &invoice.Invoice{}
	err = cbor.Unmarshal(cborBytes, invoice)
	return invoice, err
}
