package annualregistry

import (
	"ksef/internal/invoice"
	"ksef/internal/utils"
)

func (i *Invoice) AddCorrection(correction *invoice.Invoice, checksum string) error {
	correctionDataBytes, err := correction.Marshall()
	if err != nil {
		return err
	}

	correctionEntry := &Invoice{
		Correction: true,
		RefNo:      correction.Number,
		Contents:   utils.Base64ChunkedString(correctionDataBytes, 80), // let's use a chunk size of 80 characters for clarity
		// ContentChecksum: utils.Sha256Hex(correctionDataBytes),               // contents of the actual invoice data, **NOT** the XML
		Checksum:       checksum, // that one is actually the checksum of the generated XML file
		GenerationTime: correction.GenerationTime,
	}
	i.Corrections = append(i.Corrections, correctionEntry)

	return nil
}
