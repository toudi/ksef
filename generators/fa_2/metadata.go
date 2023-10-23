package fa_2

import (
	"embed"
	_ "embed"
	"ksef/metadata"
)

//go:embed "metadata.xml"
var fa_2_metadata_template embed.FS

func (fg *FA2Generator) PopulateMetadata(meta *metadata.Metadata, sourceFile string) error {
	return meta.Prepare(sourceFile, fa_2_metadata_template)
}
