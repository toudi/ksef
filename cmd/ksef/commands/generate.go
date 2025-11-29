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

type generateArgsType struct {
	Output                 string
	Delimiter              string
	GeneratorName          string
	EncodingConversionFile string
	SheetName              string
	Offline                bool
}

var generateArgs = &generateArgsType{}

func init() {
	flags := generateCommand.Flags()
	flags.StringVarP(&generateArgs.Output, "output", "o", "", "nazwa katalogu wyjściowego")
	flags.StringVarP(&generateArgs.Delimiter, "csv.delimiter", "d", ",", "łańcuch znaków rozdzielający pola (tylko dla CSV)")
	flags.StringVarP(&generateArgs.EncodingConversionFile, "csv.encoding", "e", "", "użyj pliku z konwersją znaków (tylko dla CSV)")
	flags.StringVarP(&generateArgs.SheetName, "xlsx.sheet", "s", "", "Nazwa skoroszytu do przetworzenia (tylko dla XLSX)")
	flags.StringVarP(&generateArgs.GeneratorName, "generator", "g", "fa-3_1.0", "nazwa generatora")
	flags.BoolVar(&generateArgs.Offline, "offline", false, "oznacz faktury jako generowane w trybie offline")
	flags.BoolP("mkdir", "m", false, "stwórz katalog rejestru, jeśli nie istnieje")

	RootCommand.AddCommand(generateCommand)
}

func generateRun(cmd *cobra.Command, args []string) error {
	logging.GenerateLogger.Info("generate")
	fileName := args[0]
	if fileName == "" || generateArgs.Output == "" {
		return cmd.Usage()
	}

	vip := viper.GetViper()

	// check if the registry directory does not exist:
	if _, err := os.Stat(generateArgs.Output); os.IsNotExist(err) {
		if !viper.GetBool("mkdir") {
			return errTargetDirectoryDoesNotExist
		}
		if err = os.MkdirAll(generateArgs.Output, 0775); err != nil {
			return err
		}
	}

	registry, err := registryPkg.OpenOrCreate(generateArgs.Output)
	if err != nil {
		return err
	}

	var conversionParameters inputprocessors.InputProcessorConfig
	conversionParameters.CSV.Delimiter = generateArgs.Delimiter
	conversionParameters.CSV.EncodingConversionFile = generateArgs.EncodingConversionFile
	conversionParameters.XLSX.SheetName = generateArgs.SheetName
	conversionParameters.Generator = generateArgs.GeneratorName
	conversionParameters.OfflineMode = generateArgs.Offline

	var xmlBuffer bytes.Buffer

	sei, err := sei.SEI_Init(
		runtime.GetGateway(viper.GetViper()),
		conversionParameters,
		sei.WithInvoiceReadyFunc(func(i *sei.ParsedInvoice) error {
			xmlBuffer.Reset()
			outputFile, err := os.Create(path.Join(generateArgs.Output, fmt.Sprintf("invoice-%d.xml", len(registry.Invoices))))
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
			if generateArgs.Offline {
				if certificate, err = getOfflineCertificate(
					vip,
					i.Invoice.IssuerNIP,
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
						NIP: i.Invoice.IssuerNIP,
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
