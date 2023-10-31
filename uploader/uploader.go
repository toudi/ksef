package uploader

import (
	"fmt"
	"ksef/common"
	"ksef/common/aes"
	"ksef/metadata"
	"os"
)

type Uploader struct {
	TestGateway     bool
	issuer          string
	host            string
	certificateFile string
	token           string
	cipher          *aes.Cipher
}

func (u *Uploader) Upload(sourcePath string, interactive bool) error {
	var err error

	u.host = common.KSeFHost
	u.certificateFile = common.KSeFCertificate

	if u.TestGateway {
		u.host = common.KSeFTestHost
		u.certificateFile = common.KSeFTestCertificate
	}

	if u.cipher, err = aes.CipherInit(32); err != nil {
		return fmt.Errorf("unable to init encryption cipher: %v", err)
	}

	if interactive {
		if u.issuer, err = metadata.ParseIssuerFromInvoice(sourcePath + string(os.PathSeparator) + "invoice-0.xml"); err != nil {
			return fmt.Errorf("nie rozpoznano numeru NIP: %v", err)
		}

		return u.uploadInteractive(sourcePath)
	}
	return u.uploadBatch(sourcePath)
}
