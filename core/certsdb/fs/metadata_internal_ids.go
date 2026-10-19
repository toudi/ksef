package fs

import (
	"context"
	"errors"
	"iter"
	"slices"
	"strings"

	"github.com/samber/lo"
)

func (f *FSCertsStorage) AddInternalIDsToCertificate(
	ctx context.Context,
	certId string,
	internalIds iter.Seq[string],
) error {
	for _, cert := range f.certs {
		if !(cert.UID == certId || cert.SerialNumber == certId || cert.ProfileName == certId) {
			continue
		}

		for internalId := range internalIds {
			internalIdParts := strings.Split(internalId, "-")
			if !slices.Contains(cert.NIP, internalIdParts[0]) {
				return errors.Join(ErrNIPNotOnTheList, errors.New(internalIdParts[0]))
			}
			if !slices.Contains(cert.InternalIDs, internalId) {
				cert.InternalIDs = append(cert.InternalIDs, internalId)
			}
		}
	}

	f.dirty = true

	return nil
}

func (f *FSCertsStorage) RemoveInternalIDsFromCertificate(
	ctx context.Context,
	certId string,
	internalIds iter.Seq[string],
) error {
	for _, cert := range f.certs {
		if !(cert.UID == certId || cert.SerialNumber == certId || cert.ProfileName == certId) {
			continue
		}

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
	}

	f.dirty = true

	return nil
}
