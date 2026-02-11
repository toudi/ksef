package invoice

import (
	"bytes"
	"compress/gzip"
	"errors"

	"github.com/fxamacker/cbor/v2"
)

var (
	errMarshallingInvoice = errors.New("unable to marshall invoice")
	errCompressingInvoice = errors.New("unable to compress invoice")
	errClosingGZIPWriter  = errors.New("unable to close gzip writer")
)

func (i *Invoice) Marshall() ([]byte, error) {
	// we're going to use cbor, then pass it to gzip
	// and encode the result bytes as base64 for space
	// efficiency
	invoiceBytes, err := cbor.Marshal(i)
	if err != nil {
		return nil, errors.Join(errMarshallingInvoice, err)
	}
	// once we have the cbor bytes, we can gzip them
	var compressedInvoiceBytesBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedInvoiceBytesBuffer)
	if _, err = gzipWriter.Write(invoiceBytes); err != nil {
		return nil, errors.Join(errCompressingInvoice, err)
	}
	if err = gzipWriter.Close(); err != nil {
		return nil, errors.Join(errClosingGZIPWriter, err)
	}

	return compressedInvoiceBytesBuffer.Bytes(), nil
}

// func (i *Invoice) ContentChecksum() (string, error) {
// 	contents, err := i.Marshall()
// 	if err != nil {
// 		return "", err
// 	}

// 	return utils.Sha256Hex(contents), nil
// }
