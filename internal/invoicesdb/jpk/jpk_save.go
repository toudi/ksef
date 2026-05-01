package jpk

import (
	"ksef/internal/xml"
	"os"
	"path/filepath"
)

const (
	JPKPurchase = "JPK.Ewidencja.ZakupWiersz."
	JPKIncome   = "JPK.Ewidencja.SprzedazWiersz."
)

func (j *JPK) writeToFile(root *xml.Node, outputDir string) error {
	var err error
	if err = os.MkdirAll(outputDir, 0775); err != nil {
		return err
	}
	outputFilename := "jpk-v7m.xml"
	// TODO: properly recognize if we're creating a correction. if so,
	// set the following path:
	// root.SetValue("JPK.Naglowek.CelZlozenia", "2")
	writer, err := os.Create(filepath.Join(outputDir, outputFilename))
	if err != nil {
		return err
	}
	defer writer.Close()
	return root.DumpToWriter(writer, 0)
}
