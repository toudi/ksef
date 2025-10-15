package encryption

import (
	"encoding/base64"
)

type CipherTemplateVarsType struct {
	IV            []byte
	EncryptionKey []byte
}

type CipherHTTPRequest struct {
	Key string `json:"encryptedSymmetricKey"`
	IV  string `json:"initializationVector"`
}

func (c Cipher) PrepareHTTPRequestPayload(certificateFile string) (CipherHTTPRequest, error) {
	var chr CipherHTTPRequest
	var encoder = base64.StdEncoding

	encryptedKeyBytes, err := EncryptMessageWithCertificate(certificateFile, c.Key)
	if err != nil {
		return chr, err
	}
	chr.Key = encoder.EncodeToString(encryptedKeyBytes)
	chr.IV = encoder.EncodeToString(c.IV)
	return chr, nil
}
