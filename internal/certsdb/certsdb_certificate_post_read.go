package certsdb

import (
	"crypto/x509"
	"encoding/pem"
	"ksef/internal/logging"
	"ksef/internal/runtime"
	"os"
	"path"
	"path/filepath"
	"slices"
)

func (cert *Certificate) postReadHook(environmentId string) bool {
	var updateRequired bool

	updateRequired = cert.replaceLegacyEnvironmentReferences()
	// load up only the certificates that belong to the selected environment
	cert.available = cert.EnvironmentId == environmentId

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
			return updateRequired || false
		}
		derBytes, _ := pem.Decode(pemBytes)
		certInfo, err := x509.ParseCertificate(derBytes.Bytes)
		if err != nil {
			logging.CertsDBLogger.Error("unable to prase cert info", "err", err)
			return updateRequired || false
		}
		cert.CN = &certInfo.Subject.CommonName
		return true
	}

	return updateRequired || false
}

func (cert *Certificate) replaceLegacyEnvironmentReferences() bool {
	currentEnvironmentId := cert.EnvironmentId

	for newEnvironmentId, legacyHosts := range runtime.LegacyEnvironmentHosts {
		for _, host := range legacyHosts {
			if currentEnvironmentId == host {
				// we have a match for a legacy environment ID (which was based on hostname)
				// replace EnvironmentId so that we can use the new scheme
				cert.EnvironmentId = newEnvironmentId

				// filenames were based on the legacy host
				legacyCertFilename := filepath.Join(path.Dir(certificatesDBFile), host+"-"+cert.UID+".pem")
				// rename files
				if err := os.Rename(legacyCertFilename, cert.Filename()); err != nil {
					panic(err)
				}
				if slices.ContainsFunc(cert.Usage, func(usage Usage) bool {
					return usage == UsageOffline || usage == UsageAuthentication
				}) {
					legacyPKFilename := filepath.Join(path.Dir(certificatesDBFile), host+"-"+cert.UID+"-pkey.pem")

					if err := os.Rename(legacyPKFilename, cert.PrivateKeyFilename()); err != nil {
						panic(err)
					}
				}

				return true
			}
		}
	}

	return false
}
