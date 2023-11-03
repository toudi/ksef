package api

import (
	"embed"
	"encoding/base64"
	"fmt"
	"ksef/common"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

type FileSizeAndHash struct {
	Size int
	Hash []byte
}

type metadataTemplateVarsType struct {
	Cipher           cipherTemplateVarsType
	Issuer           string
	Archive          FileSizeAndHash
	EncryptedArchive FileSizeAndHash
}

//go:embed "batch_metadata.xml"
var batchMetadataFile embed.FS

var metadataTemplateVars metadataTemplateVarsType

func (b *BatchSession) GenerateMetadata(sourcePath string) error {
	var err error

	collection, err := common.InvoiceCollection(sourcePath)
	metadataTemplateVars.Issuer = collection.Issuer
	metadataTemplateVars.Cipher = b.api.cipherTemplateVars

	if err != nil {
		return fmt.Errorf("cannot parse invoice collection: %v", err)
	}

	archive, err := Archive_init(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot create archive file: %v", err)
	}

	for _, fileName := range collection.Files {
		if _, err = archive.addFile(fileName); err != nil {
			return fmt.Errorf("unable to add file to archive: %v", err)
		}
	}

	archive.Close()

	encryptedArchive, err := archive.encrypt(b.api.cipher)
	if err != nil {
		return fmt.Errorf("cannot encrypt archive: %v", err)
	}

	if metadataTemplateVars.Cipher.EncryptionKey, err = common.EncryptMessageWithCertificate(b.api.environment.rsaPublicKey, b.api.cipher.Key); err != nil {
		return fmt.Errorf("unable to encrypt aes key: %v", err)
	}

	metadataTemplateVars.Archive.Hash = make([]byte, len(archive.hash))
	copy(metadataTemplateVars.Archive.Hash, archive.hash)
	metadataTemplateVars.Archive.Size = archive.filesize
	metadataTemplateVars.EncryptedArchive.Size = encryptedArchive.size
	metadataTemplateVars.EncryptedArchive.Hash = make([]byte, len(encryptedArchive.hash))
	copy(metadataTemplateVars.EncryptedArchive.Hash, encryptedArchive.hash)

	var funcMap = template.FuncMap{
		"base64":   base64.StdEncoding.EncodeToString,
		"filename": path.Base,
	}

	tmpl, err := template.New("batch_metadata.xml").Funcs(funcMap).ParseFS(batchMetadataFile, "batch_metadata.xml")
	if err != nil {
		return fmt.Errorf("cannot parse template: %v", err)
	}
	outputFile, err := os.Create(filepath.Join(sourcePath, "metadata.xml"))
	if err != nil {
		return fmt.Errorf("cannot create output file: %v", err)
	}

	return tmpl.Execute(outputFile, metadataTemplateVars)

}
