package metadata

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const metadataFileName = "metadata.xml"
const archiveFileName = metadataFileName + ".zip"

type Metadata struct {
	CertificateFile string
	Issuer          string
}

type MetadataTemplateVars struct {
	Cipher struct {
		IV            []byte
		EncryptionKey []byte
	}
	Archive struct {
		Size int
		Hash []byte
	}
	EncryptedArchive struct {
		Hash []byte
		Size int
	}
	Issuer string
}

func (m *Metadata) Prepare(sourcePath string, metadataTemplate fs.FS) error {
	var err error
	archive, err := Archive_init(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot create archive file: %v", err)
	}

	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot read list of files from %s: %v", sourcePath, err)
	}

	var fileName string

	for _, file := range files {
		fileName = file.Name()
		if filepath.Ext(fileName) != ".xml" || filepath.Base(fileName) == metadataFileName {
			continue
		}

		_, err = archive.addFile(fileName)
		if err != nil {
			return fmt.Errorf("cannot add %s to archive: %v", fileName, err)
		}
	}

	archive.Close()

	encryptedArchive, err := archive.encrypt()
	if err != nil {
		return fmt.Errorf("cannot encrypt archive: %v", err)
	}

	encryptedKey, err := encryptedArchive.encryptKeyWithCertificate(m.CertificateFile)
	if err != nil {
		return fmt.Errorf("cannot encrypt the key with ceritifcate file %s: %v", m.CertificateFile, err)
	}

	templateVars := &MetadataTemplateVars{
		Issuer: m.Issuer,
	}
	templateVars.Cipher.IV = make([]byte, aes.BlockSize)
	copy(templateVars.Cipher.IV, encryptedArchive.cipher.IV)
	templateVars.Cipher.EncryptionKey = make([]byte, len(encryptedKey))
	copy(templateVars.Cipher.EncryptionKey, encryptedKey)
	templateVars.Archive.Hash = make([]byte, len(archive.hash))
	copy(templateVars.Archive.Hash, archive.hash)
	templateVars.Archive.Size = archive.filesize
	templateVars.EncryptedArchive.Size = encryptedArchive.size
	templateVars.EncryptedArchive.Hash = make([]byte, len(encryptedArchive.hash))
	copy(templateVars.EncryptedArchive.Hash, encryptedArchive.hash)

	var funcMap = template.FuncMap{
		"base64":   base64.StdEncoding.EncodeToString,
		"filename": path.Base,
	}

	tmpl, err := template.New("metadata.xml").Funcs(funcMap).ParseFS(metadataTemplate, "metadata.xml")
	if err != nil {
		return fmt.Errorf("cannot parse template: %v", err)
	}
	outputFile, err := os.Create(filepath.Join(sourcePath, metadataFileName))
	if err != nil {
		return fmt.Errorf("cannot create output file: %v", err)
	}

	return tmpl.Execute(outputFile, templateVars)
}
