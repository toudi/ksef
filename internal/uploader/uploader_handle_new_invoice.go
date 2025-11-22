package uploader

import (
	"errors"
	"fmt"
	"io"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/sei"
	"os"
	"path"

	"github.com/mozillazg/go-slugify"
)

const (
	dirnameIncome   = "wystawione"
	invoiceFilename = "%04d-%s-%s.xml"
)

func (u *Uploader) handleNewInvoice(i *sei.ParsedInvoice, checksum string) error {
	ordNo := 1
	recipientName := i.Invoice.RecipientName
	refNo := i.Invoice.Number

	var destFileName = path.Join(
		u.prefixDir, i.Invoice.Issued.Format("01"), dirnameIncome,
		fmt.Sprintf(
			invoiceFilename,
			ordNo,
			slugify.Slugify(refNo),
			slugify.Slugify(recipientName),
		),
	)
	var destDir = path.Dir(destFileName)

	if err := os.MkdirAll(destDir, 0775); err != nil {
		return err
	}

	destFile, err := os.Create(destFileName)
	if err != nil {
		return errors.Join(errors.New("error opening dest file"), err)
	}

	defer destFile.Close()

	var offlineCert *certsdb.Certificate

	if i.Invoice.KSeFFlags.Offline {
		if offlineCert, err = u.GetOfflineModeCertificate(i.Invoice.IssuerNIP); err != nil {
			return errors.Join(errors.New("error selecting certificate"), err)
		}
	}

	if err = u.registry.AddInvoice(
		invoices.InvoiceMetadata{
			Metadata:      i.Invoice.Meta,
			InvoiceNumber: i.Invoice.Number,
			IssueDate:     i.Invoice.Issued.Format("2006-01-02"),
			Seller: invoices.InvoiceSubjectMetadata{
				NIP: i.Invoice.IssuerNIP,
			},
			Offline: i.Invoice.KSeFFlags.Offline,
		},
		checksum,
		offlineCert,
	); err != nil {
		return errors.Join(errors.New("error saving to registry"), err)
	}

	if _, err = io.Copy(destFile, &u.contentBuffer); err != nil {
		return errors.Join(errors.New("error copying to dest file"), err)
	}

	u.invoiceDB.Add(i.Invoice, checksum)

	return nil
}
