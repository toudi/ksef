package fa_1

import (
	"embed"
	_ "embed"
	"ksef/metadata"
)

//go:embed "metadata.xml"
var fa_1_1_metadata_template embed.FS

func (fg *FA1Generator) PopulateMetadata(meta *metadata.Metadata, sourceFile string) error {
	return meta.Prepare(sourceFile, fa_1_1_metadata_template)
}
