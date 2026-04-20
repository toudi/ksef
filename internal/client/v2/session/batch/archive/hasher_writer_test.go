package archive

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasherWriter(t *testing.T) {
	// so our input is 10 bytes and we'd like to chunk it by 4 bytes each
	inputData := bytes.NewBuffer([]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09})
	// expected sums were calculated in python like so:
	// from hashlib import sha256
	// from base64 import b64encode
	//
	// >>> b64encode(sha256(b"\x00\x01\x02\x03").digest())
	// 'BU7ewdAhH2JP7Qy8qdT5QAsOSRxDdCryxbCr6/DJkNg='
	// >>> b64encode(sha256(b"\x04\x05\x06\x07").digest())
	// 'xtRM9Bj2EOP+nh2SlP9D3vgcbNytbLsYIM/0jTqkNV0='
	// >>> b64encode(sha256(b"\x08\x09").digest())
	// 'c5B1iRAafoq4MXjn2ymXqrcnLNAtNk6OPswr7M2ktjE='

	type hashAndSize struct {
		Hash         string
		BytesWritten int
	}

	hashes := []hashAndSize{}
	var buffers []bytes.Buffer

	getWriter := func() (io.Writer, error) {
		var newBuffer bytes.Buffer
		buffers = append(buffers, newBuffer)
		return &buffers[len(buffers)-1], nil
	}
	chunkWritten := func(hash string, bytesWritten int) {
		hashes = append(hashes, hashAndSize{Hash: hash, BytesWritten: bytesWritten})
	}

	hw, err := NewHasherWriter(inputData, 4, getWriter, chunkWritten)
	require.NoError(t, err)
	bytesWritten, err := io.Copy(hw, inputData)
	_, _ = hw.Write(nil)
	require.NoError(t, err)
	require.Equal(t, int64(10), bytesWritten)

	expectedContents := [][]byte{
		{0x00, 0x01, 0x02, 0x03},
		{0x04, 0x05, 0x06, 0x07},
		{0x08, 0x09},
	}
	expectedHashes := []hashAndSize{
		{Hash: "BU7ewdAhH2JP7Qy8qdT5QAsOSRxDdCryxbCr6/DJkNg=", BytesWritten: 4},
		{Hash: "xtRM9Bj2EOP+nh2SlP9D3vgcbNytbLsYIM/0jTqkNV0=", BytesWritten: 4},
		{Hash: "c5B1iRAafoq4MXjn2ymXqrcnLNAtNk6OPswr7M2ktjE=", BytesWritten: 2},
	}
	for bufferIndex, buffer := range buffers {
		require.Equal(t, expectedContents[bufferIndex], buffer.Bytes())
		require.Equal(t, expectedHashes[bufferIndex], hashes[bufferIndex])
	}
}
