package archive

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasherWriter(t *testing.T) {
	// so our input is 10 bytes and we'd like to chunk it by 4 bytes each
	var inputData = bytes.NewBuffer([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09})
	// expected sums were calculated in python like so:
	// from hashlib import sha256
	//
	// >>> sha256(b"\x00\x01\x02\x03").hexdigest()
	// '054edec1d0211f624fed0cbca9d4f9400b0e491c43742af2c5b0abebf0c990d8'
	// >>> sha256(b"\x04\x05\x06\x07").hexdigest()
	// 'c6d44cf418f610e3fe9e1d9294ff43def81c6cdcad6cbb1820cff48d3aa4355d'
	// >>> sha256(b"\x08\x09").hexdigest()
	// '73907589101a7e8ab83178e7db2997aab7272cd02d364e8e3ecc2beccda4b631'

	type hashAndSize struct {
		Hash         string
		BytesWritten int
	}

	var hashes = []hashAndSize{}
	var buffers []bytes.Buffer

	var getWriter = func() (io.Writer, error) {
		var newBuffer bytes.Buffer
		buffers = append(buffers, newBuffer)
		return &buffers[len(buffers)-1], nil
	}
	var chunkWritten = func(hash string, bytesWritten int) {
		hashes = append(hashes, hashAndSize{Hash: hash, BytesWritten: bytesWritten})
	}

	hw, err := NewHasherWriter(inputData, 4, getWriter, chunkWritten)
	require.NoError(t, err)
	bytesWritten, err := io.Copy(hw, inputData)
	_, _ = hw.Write(nil)
	require.NoError(t, err)
	require.Equal(t, int64(10), bytesWritten)

	var expectedContents = [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		{0x04, 0x05, 0x06, 0x07},
		{0x08, 0x09},
	}
	var expectedHashes = []hashAndSize{
		{Hash: "054edec1d0211f624fed0cbca9d4f9400b0e491c43742af2c5b0abebf0c990d8", BytesWritten: 4},
		{Hash: "c6d44cf418f610e3fe9e1d9294ff43def81c6cdcad6cbb1820cff48d3aa4355d", BytesWritten: 4},
		{Hash: "73907589101a7e8ab83178e7db2997aab7272cd02d364e8e3ecc2beccda4b631", BytesWritten: 2},
	}
	for bufferIndex, buffer := range buffers {
		require.Equal(t, expectedContents[bufferIndex], buffer.Bytes())
		require.Equal(t, expectedHashes[bufferIndex], hashes[bufferIndex])
	}
}
