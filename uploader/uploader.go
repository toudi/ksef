package uploader

import (
	"fmt"
	"ksef/common"
	"ksef/metadata"
	"os"
)

type Uploader struct {
	TestGateway     bool
	issuer          string
	host            string
	certificateFile string
	token           string
}

func (u *Uploader) Upload(sourcePath string, interactive bool) error {
	var err error

	u.host = common.KSeFHost
	u.certificateFile = common.KSeFCertificate

	if u.TestGateway {
		u.host = common.KSeFTestHost
		u.certificateFile = common.KSeFTestCertificate
	}

	if interactive {
		if u.issuer, err = metadata.ParseIssuerFromInvoice(sourcePath + string(os.PathSeparator) + "invoice-0.xml"); err != nil {
			return fmt.Errorf("nie rozpoznano numeru NIP: %v", err)
		}

		return u.uploadInteractive(sourcePath)
	}
	return u.uploadBatch(sourcePath)
}
