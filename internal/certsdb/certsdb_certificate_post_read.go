package certsdb

import (
	"crypto/x509"
	"encoding/pem"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"os"
)

func (cert *Certificate) postReadHook(environment runtime.Gateway) bool {
	// load up only the certificates that belong to the selected environment
	cert.available = cert.Environment == environment

	// check if the loaded NIP's are a slice of strings or just a string
	if nipString, ok := cert.NIPRaw.(string); ok {
		cert.NIP = []string{nipString}
	}
	if nipSliceRaw, ok := cert.NIPRaw.([]any); ok {
		for _, nipAny := range nipSliceRaw {
			if nipString, ok := nipAny.(string); ok {
				cert.NIP = append(cert.NIP, nipString)
			}
		}
	}

	if cert.CN == nil {
		// read CN from the cert file itself
		pemBytes, err := os.ReadFile(cert.Filename())
		if err != nil {
			return false
		}
		derBytes, _ := pem.Decode(pemBytes)
		certInfo, err := x509.ParseCertificate(derBytes.Bytes)
		if err != nil {
			logging.CertsDBLogger.Error("unable to prase cert info", "err", err)
			return false
		}
		cert.CN = &certInfo.Subject.CommonName
		return true
	}

	return false
}
