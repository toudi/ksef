package uploader

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

func unmarshallInvoice(content string) (*invoice.Invoice, error) {
	// let's do the exact opposite of what marshall does.
	// this means:
	// 1. treat input as gzip content, base64 encoded
	gzipBytes, err := base64.StdEncoding.DecodeString(content)
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
	var invoice = &invoice.Invoice{}
	err = cbor.Unmarshal(cborBytes, invoice)
	return invoice, err
}

func (i Invoice) SourceDocument() (*invoice.Invoice, error) {
	return unmarshallInvoice(i.Contents)
}
