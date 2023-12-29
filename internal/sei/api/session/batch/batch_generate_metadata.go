package batch

import (
	"embed"
	"encoding/base64"
	"fmt"
	encryptionPkg "ksef/internal/encryption"
	"ksef/internal/invoice"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
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
	Cipher           encryptionPkg.CipherTemplateVarsType
	Issuer           string
	Archive          FileSizeAndHash
	EncryptedArchive FileSizeAndHash
}

//go:embed "batch_metadata.xml"
var batchMetadataFile embed.FS

var metadataTemplateVars metadataTemplateVarsType

func (b *BatchSession) GenerateMetadata(sourcePath string) error {
	var err error

	// we try to load the registry. If it does not exists, it's not a fatal problem since
	// we will receive a new registry as the return value and it will still conform to the
	// interface, it will simply report each invoice as not sent.
	var registryFile string = path.Join(sourcePath, "registry.yaml")
	registry, err := registryPkg.OpenOrCreate(registryFile)
	if err != nil {
		return fmt.Errorf("cannot open registry: %v", err)
	}

	collection, err := invoice.InvoiceCollection(sourcePath, registry)
	if err == invoice.ErrAlreadySynced {
		logging.UploadLogger.Info("no invoices left to send")
		return nil
	}

	if err != nil {
		return fmt.Errorf("cannot parse invoice collection: %v", err)
	}

	encryption, err := b.apiClient.Encryption()
	if err != nil {
		return fmt.Errorf("cannot instantiate cipher: %v", err)
	}

	metadataTemplateVars.Issuer = collection.Issuer
	metadataTemplateVars.Cipher = encryption.CipherTemplateVars

	archive, err := Archive_init(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot create archive file: %v", err)
	}

	for _, file := range collection.Files {
		if _, err = archive.addFile(file.Filename); err != nil {
			return fmt.Errorf("unable to add file to archive: %v", err)
		}
	}

	archive.Close()

	encryptedArchive, err := archive.encrypt(encryption.Cipher)
	if err != nil {
		return fmt.Errorf("cannot encrypt archive: %v", err)
	}

	if metadataTemplateVars.Cipher.EncryptionKey, err = encryptionPkg.EncryptMessageWithCertificate(b.apiClient.Environment.RsaPublicKey, encryption.Cipher.Key); err != nil {
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

	tmpl, err := template.New("batch_metadata.xml").
		Funcs(funcMap).
		ParseFS(batchMetadataFile, "batch_metadata.xml")
	if err != nil {
		return fmt.Errorf("cannot parse template: %v", err)
	}
	outputFile, err := os.Create(filepath.Join(sourcePath, "metadata.xml"))
	if err != nil {
		return fmt.Errorf("cannot create output file: %v", err)
	}

	if err = tmpl.Execute(outputFile, metadataTemplateVars); err != nil {
		return fmt.Errorf("cannot generate metadata file: %v", err)
	}

	if err = registry.Save(registryFile); err != nil {
		return fmt.Errorf("cannot save registry: %v", err)
	}

	return nil
}
