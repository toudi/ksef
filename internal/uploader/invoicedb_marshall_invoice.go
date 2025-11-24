package uploader

import (
	"bytes"
	"compress/gzip"
	"errors"
	"ksef/internal/invoice"
	"ksef/internal/utils"

	"github.com/fxamacker/cbor/v2"
)

var (
	errMarshallingInvoice = errors.New("unable to marshall invoice")
	errCompressingInvoice = errors.New("unable to compress invoice")
	errClosingGZIPWriter  = errors.New("unable to close gzip writer")
)

func marshallInvoice(i *invoice.Invoice) (string, error) {
	// we're going to use cbor, then pass it to gzip
	// and encode the result bytes as base64 for space
	// efficiency
	invoiceBytes, err := cbor.Marshal(i)
	if err != nil {
		return "", errors.Join(errMarshallingInvoice, err)
	}
	// once we have the cbor bytes, we can gzip them
	var compressedInvoiceBytesBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedInvoiceBytesBuffer)
	if _, err = gzipWriter.Write(invoiceBytes); err != nil {
		return "", errors.Join(errCompressingInvoice, err)
	}
	if err = gzipWriter.Close(); err != nil {
		return "", errors.Join(errClosingGZIPWriter, err)
	}
	// finally - read bytes from compressed buffer and write them to base64
	// let's use a chunk size of 80 for clarity
	return utils.Base64ChunkedString(compressedInvoiceBytesBuffer.Bytes(), 80), nil
}
