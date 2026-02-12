package certsdb

import (
	"errors"
	"iter"
	"slices"
	"strings"

	"github.com/samber/lo"
)

var errNIPNotOnTheList = errors.New("NIP is not on the list of this certificate's NIP's")

func (cdb *CertificatesDB) AddInternalIDsToCertificate(
	internalIds iter.Seq[string],
	certificateId string,
) error {
	return cdb.ForEach(
		func(cert Certificate) bool {
			return cert.UID == certificateId || cert.SerialNumber == certificateId || cert.ProfileName == certificateId
		},
		func(cert *Certificate) error {
			for internalId := range internalIds {
				internalIdParts := strings.Split(internalId, "-")
				if !slices.Contains(cert.NIP, internalIdParts[0]) {
					return errors.Join(errNIPNotOnTheList, errors.New(internalIdParts[0]))
				}
				if !slices.Contains(cert.InternalIDs, internalId) {
					cert.InternalIDs = append(cert.InternalIDs, internalId)
				}
			}
			return nil
		},
	)
}

func (cdb *CertificatesDB) RemoveInternalIdsFromCertificate(
	internalIds iter.Seq[string],
	certificateId string,
) error {
	return cdb.ForEach(
		func(cert Certificate) bool {
			return cert.UID == certificateId || cert.SerialNumber == certificateId || cert.ProfileName == certificateId
		},
		func(cert *Certificate) error {
			// create a temporary set so that we don't reallocate memory over and over
			// each time we'd want to rebuild the final list
			tmpInternalIds := make(map[string]struct{})
			for _, internalId := range cert.InternalIDs {
				tmpInternalIds[internalId] = struct{}{}
			}

			for internalId := range internalIds {
				delete(tmpInternalIds, internalId)
			}

			cert.InternalIDs = lo.Keys(tmpInternalIds)
			return nil
		},
	)
}
