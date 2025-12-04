package uploader

import (
	"errors"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/registry"
	"ksef/internal/sei"
	"ksef/internal/utils"
	"os"
	"path"
	"time"
)

const (
	registriesDir = "wysylki"
)

var (
	errUnableToRenderXML   = errors.New("unable to render XML to temporary buffer")
	errUnableToHash        = errors.New("unable to hash temporary invoice")
	ErrAutoCorrectDisabled = errors.New("auto-correct is disabled")
)

func (u *Uploader) InvoiceReady(i *sei.ParsedInvoice) error {
	var err error
	u.contentBuffer.Reset()
	// the prefix will be in the following format:
	// data/1111111111/2025/
	//                      01/
	//                        registry.yaml
	//                        wystawione/
	//                        otrzymane/
	//                      02/
	//                      wysylki/01/
	//                      wysylki/01/2025-01-02T03:04:05/
	//                      invoices.yaml
	fmt.Printf("invoice: %+v\n", i)
	invoice := i.Invoice
	prefix := path.Join(u.dataDir, invoice.IssuerNIP, invoice.Issued.Format("2006"))
	if err = os.MkdirAll(prefix, 0775); err != nil {
		return errors.Join(errors.New("error mkdirall prefix"), err)
	}
	u.prefixDir = prefix
	// check if we have opened annual invoices db to check for corrections
	if u.invoiceDB == nil {
		if u.invoiceDB, err = InvoiceDB_OpenOrCreate(u.vip, prefix); err != nil {
			return errors.Join(err, errors.New("invoicedb::open-or-create"))
		}
	}
	// first of all, check if we have an opened registry
	if u.registry == nil {
		var now = time.Now().Local()
		var registryName = now.Format("2006-01-02T15:04:05")
		var registryDir = path.Join(prefix, registriesDir, now.Format("01"), registryName)
		u.registry = registry.NewRegistryInDir(registryDir)
	}

	// it will default to time.Zero, which means that when passed to ToXML() function
	// it won't get overriden, unless ..
	var lastKnownGenerationTimestamp time.Time
	var _invoice *Invoice
	if _invoice = u.invoiceDB.GetByRefNo(invoice.Number); _invoice != nil {
		// unless we already have this invoice in our system.
		// the purpose of this trick is to always generate the same XML content so that
		// if we're parsing the same invoice (which has already been sent to KSeF) so that
		// we can detect this.
		lastKnownGenerationTimestamp = _invoice.GenerationTime
	}

	if err = i.ToXML(lastKnownGenerationTimestamp, &u.contentBuffer); err != nil {
		return errors.Join(errUnableToRenderXML, err)
	}

	// perfect. now let's calculate checksum
	checksum := utils.Sha256Hex(u.contentBuffer.Bytes())

	// and now we can finally detect if this is potentially a correction candidate or
	// simply we've already processed this.
	if _invoice != nil && checksum == _invoice.Checksum {
		logging.GenerateLogger.Debug("faktura już wysłana do KSeF. no-op.", "numer faktury w KSeF", _invoice.KSeFRefNo)
		return nil
	}

	// is this a new invoice ? or..
	// the invoice wasn't sent yet. no problem - just override it
	if _invoice == nil || _invoice.KSeFRefNo == "" {
		return u.handleNewInvoice(i, checksum)
	}

	// if not, it must be a correction
	if !u.autoCorrect {
		return ErrAutoCorrectDisabled
	}
	return u.handleCorrection(i, _invoice)

	// // now we can deal with the invoice
	// if _invoice := u.invoiceDB.GetByRefNo(invoice.Number); invoice != nil {
	// 	// either the invoice was already sent or we need to issue a correction.
	// 	// let's check both cases
	// 	srcChecksum := _invoice.Checksum
	// 	// now we can compare the original checksum with the current checksum,
	// 	// based on generation time
	// 	origTimestamp := invoice.GenerationTime
	// 	invoice.GenerationTime = invoice.GenerationTime
	// 	currentChecksum, err := u.generator.InvoiceChecksum(invoice)
	// 	if err != nil {
	// 		return errors.Join(err, errors.New("invoice checksum (1)"))
	// 	}
	// 	// now let's check the checksum based on persisted generation time.
	// 	if currentChecksum != srcChecksum {
	// 		// we have a correction to deal with
	// 		// replace generation time back with the original
	// 		invoice.GenerationTime = origTimestamp
	// 	} else {
	// 		// the checksums match - all is fine. the invoice was already processed.
	// 		logging.GenerateLogger.Debug("faktura przetworzona podczas wcześniejszych wywołań - no-op.")
	// 	}
	// 	return nil
	// }

	// // if we are here then the invoice was never processed
	// currentChecksum, err := u.generator.InvoiceChecksum(invoice)
	// if err != nil {
	// 	return errors.Join(err, errors.New("invoice checksum (2)"))
	// }
	// offlineCert, err := u.generator.GetOfflineModeCertificate(invoice)
	// if err != nil {
	// 	return errors.Join(err, errors.New("get offline certificate"))
	// }

	// return u.registry.AddInvoice(
	// 	invoices.InvoiceMetadata{
	// 		Metadata:      invoice.Meta,
	// 		InvoiceNumber: invoice.Number,
	// 		IssueDate:     invoice.Issued.Format("2006-01-02"),
	// 		Seller: invoices.InvoiceSubjectMetadata{
	// 			NIP: invoice.SellerNIP,
	// 		},
	// 		Offline: invoice.KSeFFlags.Offline,
	// 	},
	// 	currentChecksum,
	// 	offlineCert,
	// )
}
