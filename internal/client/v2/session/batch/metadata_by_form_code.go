package batch

import (
	"ksef/internal/client/v2/session/batch/archive"
	"ksef/internal/client/v2/session/types"
	"ksef/internal/encryption"
)

const (
	maxArchiveSize     int64 = 5_368_709_120 // 5 GiB
	maxArchivePartSize int   = 104_857_600   // 100 MiB
)

func (b *Session) generateMetadataByFormCode(
	formCode types.InvoiceFormCode,
	files []types.Invoice,
	cipher *encryption.Cipher,
) (*batchSessionInitRequest, error) {
	var err error
	batchMetadataRequest := &batchSessionInitRequest{
		FormCode: formCode,
	}

	_archive, err := archive.New(b.workDir, maxArchiveSize)
	if err != nil {
		return nil, err
	}

	batchPayload := &BatchSessionPayload{
		Archive: _archive,
	}

	for _, file := range files {
		if err = _archive.AddFile(file.Filename); err != nil {
			if err == archive.ErrExeedsMaxSize {
				// this is not a fatal error. we can simply rerun the program and it will pick up the rest of the files
				// also.. it's highly unlikely that anybody using this program will actually have so many invoices that
				// they would exceed the limits
				break
			}
		}
		batchPayload.InvoiceChecksums = append(batchPayload.InvoiceChecksums, file.Checksum)
		// yeah, that's a bit of a problem with using batch sessions.
		// I mean theoretically if we'd like to be super pure about it,
		// we should split the payload into separate sessions - these marked
		// with offline mode and this unmarked.
		// as it is the case with multiple modules in this project -
		// - desperate times call for desperate measures.
		if file.Offline {
			batchMetadataRequest.Offline = true
		}
	}

	if err = _archive.Split(maxArchivePartSize); err != nil {
		return nil, err
	}

	archiveMeta, err := _archive.Metadata()
	if err != nil {
		return nil, err
	}
	batchMetadataRequest.BatchFile.FileSize = uint64(archiveMeta.FileSize)
	batchMetadataRequest.BatchFile.FileHash = archiveMeta.Hash

	for partNo, part := range _archive.Parts {
		if err = part.Encrypt(cipher); err != nil {
			return nil, err
		}
		batchMetadataRequest.BatchFile.FileParts = append(
			batchMetadataRequest.BatchFile.FileParts,
			batchArchivePart{
				OrdNo:    uint32(partNo + 1),
				FileSize: uint64(part.FileSize),
				FileHash: part.Hash,
			},
		)
	}

	batchMetadataRequest.Payload = batchPayload

	return batchMetadataRequest, nil
}
