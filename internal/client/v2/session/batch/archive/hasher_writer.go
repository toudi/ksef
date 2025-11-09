package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
)

type getWriterFunc = func() (io.Writer, error)
type chunkWrittenFunc = func(hash string, bytesWritten int)

// this type is a convenience wrapper so that we can both add files to archive
// and calculate their checksum at the same time.
type HasherWriter struct {
	sizeLimit    int
	bytesWritten int
	hasher       hash.Hash
	// whenever HasherWriter recognizes that it's current writer is full,
	// it will call getWriter() to obtain a new writer.
	getWriter    getWriterFunc
	chunkWritten chunkWrittenFunc
	targetBuffer io.Writer
}

func NewHasherWriter(input io.Reader, sizeLimit int, getWriter getWriterFunc, chunkWritten chunkWrittenFunc) (*HasherWriter, error) {
	hw := &HasherWriter{
		getWriter:    getWriter,
		chunkWritten: chunkWritten,
		sizeLimit:    sizeLimit,
		hasher:       sha256.New(),
	}

	var err error
	if hw.targetBuffer, err = getWriter(); err != nil {
		return nil, err
	}

	return hw, nil
}

func (h *HasherWriter) sendHash() {
	hash := hex.EncodeToString(h.hasher.Sum(nil))
	h.chunkWritten(hash, h.bytesWritten)
	h.hasher.Reset()
}

func (h *HasherWriter) Write(p []byte) (int, error) {
	var bytesRemaining = len(p)
	// offset in p from which we're starting the read
	var offset int = 0
	var err error

	// it is a special case when we're making sure that the last hash will not be lost
	if p == nil {
		h.sendHash()
		return 0, nil
	}

	// let's loop until there are no more bytes to be written and hashed.
	var replaceWriter = false
	for bytesRemaining > 0 {
		if replaceWriter {
			if h.targetBuffer, err = h.getWriter(); err != nil {
				return -1, err
			}
			h.bytesWritten = 0
		}
		availableSpace := h.sizeLimit - h.bytesWritten
		chunkToBeWritten := p[offset : offset+min(bytesRemaining, availableSpace)]

		replaceWriter = availableSpace < bytesRemaining

		if _, err = h.hasher.Write(chunkToBeWritten); err != nil {
			return -1, err
		}

		var written int
		written, err = h.targetBuffer.Write(chunkToBeWritten)

		h.bytesWritten += written
		bytesRemaining -= written
		offset += written

		// if we've written less bytes than there's a limit for, we
		// have to exchange the writer and write the remaining part
		if replaceWriter {
			h.sendHash()
		}
	}
	return len(p), err
}
