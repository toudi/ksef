package commands

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/pdf"
	"ksef/internal/registry"
	"ksef/internal/utils"
	"os"
)

type renderPDFCommand struct {
	Command
}
type renderPDFArgsType struct {
	path           string
	resolverConfig utils.FilepathResolverConfig
	invoice        string
}

var RenderPDFCommand *renderPDFCommand
var renderPDFArgs renderPDFArgsType

func init() {
	RenderPDFCommand = &renderPDFCommand{
		Command: Command{
			Name:        "render-pdf",
			FlagSet:     flag.NewFlagSet("render-pdf", flag.ExitOnError),
			Description: "drukuje PDF dla wskazanej faktury używając lokalnego szablonu",
			Run:         renderPDFRun,
		},
	}

	RenderPDFCommand.FlagSet.StringVar(
		&renderPDFArgs.path,
		"p",
		"",
		"ścieżka do pliku rejestru",
	)
	RenderPDFCommand.FlagSet.StringVar(
		&renderPDFArgs.resolverConfig.Path,
		"o",
		"",
		"ścieżka do zapisu PDF (domyślnie katalog pliku statusu + {nrRef}.pdf)",
	)
	RenderPDFCommand.FlagSet.BoolVar(
		&renderPDFArgs.resolverConfig.Mkdir,
		"m",
		false,
		"stwórz katalog, jeśli wskazany do zapisu nie istnieje",
	)
	RenderPDFCommand.FlagSet.StringVar(
		&renderPDFArgs.invoice,
		"i",
		"",
		"plik XML do wizualizacji",
	)

	registerCommand(&RenderPDFCommand.Command)
}

func renderPDFRun(c *Command) error {
	if renderPDFArgs.path == "" || renderPDFArgs.invoice == "" {
		RenderPDFCommand.FlagSet.Usage()
		return nil
	}

	registry, err := registry.LoadRegistry(renderPDFArgs.path)
	if err != nil {
		return fmt.Errorf("unable to load registry from file: %v", err)
	}

	if registry.Environment == "" {
		return fmt.Errorf("file deserialized correctly, but environment is empty")
	}

	fileContent, err := os.ReadFile(renderPDFArgs.invoice)
	if err != nil {
		return fmt.Errorf("nie udało się odczytać pliku źródłowego")
	}

	hasher := sha256.New()
	hasher.Write(fileContent)
	fileChecksum := hex.EncodeToString(hasher.Sum(nil))
	fileBase64 := base64.StdEncoding.EncodeToString(fileContent)

	logging.PDFRendererLogger.Debug("calculated checksum", "checksum", fileChecksum)

	invoiceMeta, err := registry.GetInvoiceByChecksum(fileChecksum)
	if err != nil {
		return err
	}
	if invoiceMeta.Checksum != fileChecksum {
		return fmt.Errorf("nie udało się znaleźć faktury na podstawie kryteriów wejściowych")
	}

	renderPDFArgs.resolverConfig.DefaultFilename = fmt.Sprintf(
		"%s.pdf",
		invoiceMeta.SEIReferenceNumber,
	)

	outputPath, err := utils.ResolveFilepath(renderPDFArgs.resolverConfig)
	if err == utils.ErrDoesNotExistAndMkdirNotSpecified {
		return fmt.Errorf("wskazany katalog nie istnieje a nie użyłeś opcji `-m`")
	}
	if err != nil {
		return fmt.Errorf("błąd tworzenia katalogu wyjściowego: %v", err)
	}

	printingEngine, err := pdf.GetLocalPrintingEngine()
	if err != nil {
		return fmt.Errorf("nie udało się zainicjować silnika drukowania: %v", err)
	}
	return printingEngine.Print(fileBase64, invoiceMeta, outputPath)
}
