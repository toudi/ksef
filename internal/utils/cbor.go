package utils

import (
	"bytes"
	"compress/gzip"
	"errors"

	"github.com/fxamacker/cbor/v2"
)

var (
	errMarshalling = errors.New("unable to marshall")
	errCompressing = errors.New("unable to compress")
	errClosingGZip = errors.New("unable to close gzip writer")
)

func Base64GZippedCBor(input any, width int) (string, error) {
	// we're going to use cbor, then pass it to gzip
	// and encode the result bytes as base64 for space
	// efficiency
	inputBytes, err := cbor.Marshal(input)
	if err != nil {
		return "", errors.Join(errMarshalling, err)
	}
	// once we have the cbor bytes, we can gzip them
	var compressedInputBytesBuffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedInputBytesBuffer)
	if _, err = gzipWriter.Write(inputBytes); err != nil {
		return "", errors.Join(errCompressing, err)
	}
	if err = gzipWriter.Close(); err != nil {
		return "", errors.Join(errClosingGZip, err)
	}
	// finally - read bytes from compressed buffer and write them to base64
	// let's use a chunk size of 80 for clarity
	return Base64ChunkedString(compressedInputBytesBuffer.Bytes(), 80), nil
}
