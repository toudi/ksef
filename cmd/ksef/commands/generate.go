package commands

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"ksef/internal/certsdb"
	"ksef/internal/client/v2/types/invoices"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/runtime"
	"ksef/internal/sei"
	inputprocessors "ksef/internal/sei/input_processors"
	"ksef/internal/utils"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var errTargetDirectoryDoesNotExist = errors.New("katalog docelowy nie istnieje. użyj flagi --mkdir / -m jeśli chcesz go stworzyć")

var generateCommand = &cobra.Command{
	Use:   "generate [input]",
	Short: "konwertuje plik CSV/YAML/XLSX do pliku KSeF (XML) i tworzy rejestr gotowy do wysłania",
	RunE:  generateRun,
	Args:  cobra.ExactArgs(1),
}

func init() {
	flags := generateCommand.Flags()
	inputprocessors.GeneratorFlags(flags)
	flags.StringP("output", "o", "", "nazwa katalogu wyjściowego")
	flags.BoolP("mkdir", "m", false, "stwórz katalog rejestru, jeśli nie istnieje")

	RootCommand.AddCommand(generateCommand)
}

func generateRun(cmd *cobra.Command, args []string) error {
	vip := viper.GetViper()
	output := vip.GetString("output")
	offlineMode := vip.GetBool("offline")

	logging.GenerateLogger.Info("generate")
	fileName := args[0]
	if fileName == "" || output == "" {
		return cmd.Usage()
	}

	// check if the registry directory does not exist:
	if _, err := os.Stat(output); os.IsNotExist(err) {
		if !viper.GetBool("mkdir") {
			return errTargetDirectoryDoesNotExist
		}
		if err = os.MkdirAll(output, 0775); err != nil {
			return err
		}
	}

	registry, err := registryPkg.OpenOrCreate(output)
	if err != nil {
		return err
	}

	var xmlBuffer bytes.Buffer

	sei, err := sei.SEI_Init(
		vip,

		sei.WithInvoiceReadyFunc(func(i *sei.ParsedInvoice) error {
			xmlBuffer.Reset()
			outputFile, err := os.Create(path.Join(output, fmt.Sprintf("invoice-%d.xml", len(registry.Invoices))))
			if err != nil {
				return err
			}
			defer outputFile.Close()
			if err = i.ToXML(time.Time{}, &xmlBuffer); err != nil {
				return err
			}
			checksum := utils.Sha256Hex(xmlBuffer.Bytes())
			if _, err = io.Copy(outputFile, &xmlBuffer); err != nil {
				return err
			}
			var certificate *certsdb.Certificate
			if offlineMode {
				if certificate, err = getOfflineCertificate(
					vip,
					i.Invoice.Issuer.NIP,
				); err != nil {
					return err
				}
			}
			return registry.AddInvoice(
				invoices.InvoiceMetadata{
					Metadata:      i.Invoice.Meta,
					InvoiceNumber: i.Invoice.Number,
					IssueDate:     i.Invoice.Issued.Format("2006-01-02"),
					Seller: invoices.InvoiceSubjectMetadata{
						NIP: i.Invoice.Issuer.NIP,
					},
					Offline: i.Invoice.KSeFFlags.Offline,
				},
				checksum,
				certificate,
			)
		}),
	)
	if err != nil {
		return err
	}
	if err = sei.ProcessSourceFile(fileName); err != nil {
		return fmt.Errorf("error calling processSourceFile: %v", err)
	}

	return registry.Save("")
}

var certsDB *certsdb.CertificatesDB

func getOfflineCertificate(vip *viper.Viper, nip string) (*certsdb.Certificate, error) {
	var err error

	if certsDB == nil {
		certsDB, err = certsdb.OpenOrCreate(runtime.GetGateway(vip))
		if err != nil {
			return nil, err
		}
	}

	cert, err := certsDB.GetByUsage(
		certsdb.UsageOffline,
		nip,
	)
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
