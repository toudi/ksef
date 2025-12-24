package invoicesdb

import (
	"errors"
	"io/fs"
	"ksef/internal/invoice"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/logging"
	"ksef/internal/sei"
	"ksef/internal/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	errUnableToWalkDir        = errors.New("unable to list directory")
	errUnableToParseInvoice   = errors.New("unable to parse invoice")
	errUnableToParseIssueDate = errors.New("unable to parse issue date")
	errUnableToCopyFile       = errors.New("unable to copy invoice file")
	errUnableToCreateDir      = errors.New("unable to create invoice directory")
)

func (i *InvoicesDB) importXMLInvoices(
	fileName string,
) (err error) {
	filesToImport := []string{fileName}
	// let's check if it is a single file or a wildcard (*.xml)
	fileBase := strings.ToLower(fileName)
	if fileBase == "*.xml" {
		// we need to iterate through directory, find all xml files and add them to the slice
		dirName := filepath.Dir(fileBase)
		filesToImport = []string{}

		if err = filepath.WalkDir(dirName, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if strings.ToLower(filepath.Ext(path)) == ".xml" {
				logging.InvoicesDBLogger.Debug("dodaję plik faktury do kolejki", "plik", path)
				filesToImport = append(filesToImport, path)
			}
			return nil
		}); err != nil {
			return errors.Join(errUnableToWalkDir, err)
		}
	}

	// perfect. now that we have the list of source xml files, let's go through them one by one,
	// parse the xml invoices as they go and add them to registry.
	for _, srcInvoiceFile := range filesToImport {
		invoiceMeta, checksum, err := monthlyregistry.ParseInvoice(srcInvoiceFile)
		if err != nil {
			return errors.Join(errUnableToParseInvoice, err)
		}

		issueDateString := invoiceMeta.Issued
		issueDate, err := time.Parse(time.DateOnly, issueDateString)
		if err != nil {
			return errors.Join(errUnableToParseIssueDate, err)
		}

		parsedInvoice := &sei.ParsedInvoice{
			Invoice: &invoice.Invoice{
				GenerationTime: invoiceMeta.GeneratedTime,
				Number:         invoiceMeta.InvoiceNumber,
				Issued:         issueDate,
				Issuer: invoice.Issuer{
					NIP: invoiceMeta.Issuer,
				},
				KSeFFlags: &invoice.KSeFFlags{
					Offline: i.importCfg.Offline,
				},
			},
		}

		annualRegistry, err := i.getAnnualRegistryForInvoice(
			parsedInvoice,
		)
		if err != nil {
			return err
		}

		if annualRegistry.GetByChecksum(checksum) != nil {
			logging.InvoicesDBLogger.Debug(
				"ta faktura już istnieje w bazie. no-op",
				"numer faktury", invoiceMeta.InvoiceNumber,
				"suma kontrolna", checksum,
			)
			continue
		}

		if err = annualRegistry.AddInvoice(
			parsedInvoice,
			checksum,
			false,
		); err != nil {
			return err
		}

		monthlyRegistry, err := i.getMonthlyRegistryForInvoice(
			parsedInvoice,
		)
		if err != nil {
			return err
		}

		fileName := monthlyRegistry.GetDestFileName(parsedInvoice, monthlyregistry.InvoiceTypeIssued)

		if err = os.MkdirAll(filepath.Dir(fileName), 0775); err != nil {
			return errors.Join(errUnableToCreateDir, err)
		}

		if err = utils.CopyFile(srcInvoiceFile, fileName); err != nil {
			return errors.Join(errUnableToCopyFile, err)
		}

		if err = monthlyRegistry.AddInvoice(
			parsedInvoice,
			monthlyregistry.InvoiceTypeIssued,
			checksum,
		); err != nil {
			return err
		}

		regInvoice := monthlyRegistry.GetInvoiceByChecksum(
			checksum,
		)

		i.newInvoices = append(
			i.newInvoices,
			&NewInvoice{
				registry: monthlyRegistry,
				invoice:  regInvoice,
			},
		)
	}
	return nil
}
