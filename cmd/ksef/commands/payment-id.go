package commands

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"ksef/internal/sei/api/session/interactive"
	"os"

	"gopkg.in/yaml.v3"
)

type generatePaymentIdCommand struct {
	Command
}
type generatePaymentIdArgsType struct {
	output      string
	issuerToken string
	path        string
	yaml        bool
	json        bool
}

var GeneratePaymentIdCommand *generatePaymentIdCommand
var generatePaymentIdArgs generatePaymentIdArgsType

func init() {
	GeneratePaymentIdCommand = &generatePaymentIdCommand{
		Command: Command{
			Name:        "payment-id",
			FlagSet:     flag.NewFlagSet("payment-id", flag.ExitOnError),
			Description: "generuje identyfikator płatności dla wybranych fatur",
			Run:         generatePaymentIdRun,
			Args:        generatePaymentIdArgs,
		},
	}

	GeneratePaymentIdCommand.FlagSet.StringVar(&generatePaymentIdArgs.path, "p", "", "ścieżka do pliku rejestru")
	GeneratePaymentIdCommand.FlagSet.StringVar(&generatePaymentIdArgs.output, "o", "", "Plik do zapisania wyjścia")
	GeneratePaymentIdCommand.FlagSet.BoolVar(&generatePaymentIdArgs.yaml, "yaml", false, "Użyj formatu YAML do zapisania wyjścia")
	GeneratePaymentIdCommand.FlagSet.BoolVar(&generatePaymentIdArgs.json, "json", false, "Użyj formatu JSON do zapisania wyjścia")
	GeneratePaymentIdCommand.FlagSet.StringVar(&generatePaymentIdArgs.issuerToken, "token", "", "Token sesji interaktywnej lub nazwa zmiennej środowiskowej która go zawiera")

	registerCommand(&GeneratePaymentIdCommand.Command)
}

func generatePaymentIdRun(c *Command) error {
	invoiceIds := GeneratePaymentIdCommand.FlagSet.Args()
	if len(invoiceIds) == 0 {
		GeneratePaymentIdCommand.FlagSet.Usage()
		return nil
	}

	_registry, err := registry.LoadRegistry(generatePaymentIdArgs.path)
	if err != nil {
		return fmt.Errorf("unable to load registry from file: %v", err)
	}

	if _registry.Environment == "" {
		return fmt.Errorf("file deserialized correctly, but environment is empty")
	}

	// let's check which mode are we on.
	if len(invoiceIds) == 1 {
		// this is verification mode
		invoiceIds, err := _registry.GetInvoiceIdsForPaymentId(invoiceIds[0])
		if err == registry.ErrPaymentIdNotFound {
			return fmt.Errorf("nie znaleziono identyfikatora płatności o podanym ID")
		}
		var output io.WriteCloser

		output = os.Stdout

		if generatePaymentIdArgs.output != "" {
			output, err = os.Create(generatePaymentIdArgs.output)
			if err != nil {
				return fmt.Errorf("błąd otwierania pliku wyjścia: %v", err)
			}
			defer output.Close()
		}

		if generatePaymentIdArgs.yaml {
			return yaml.NewEncoder(output).Encode(invoiceIds)
		} else if generatePaymentIdArgs.json {
			return json.NewEncoder(output).Encode(invoiceIds)
		} else {
			writer := csv.NewWriter(output)
			_ = writer.Write([]string{"invoiceRefNo", "ksefInvoiceRefNo"})
			for _, ids := range invoiceIds {
				_ = writer.Write([]string{ids.ReferenceNumber, ids.SEIReferenceNumber})
			}
			writer.Flush()
		}
		return nil
	}

	// we are in the generation mode.
	if len(invoiceIds) < 2 {
		return fmt.Errorf("stworzenie identyfikatora płatności wymaga co najmniej dwóch numerów faktur")
	}

	seiRefNumbers, err := _registry.GetSEIRefNoFromArray(invoiceIds)

	if err != nil {
		return fmt.Errorf("nie udało się odnaleźć numerów referencyjnych: %v", err)
	}

	gateway, err := client.APIClient_Init(_registry.Environment)
	if err != nil {
		return fmt.Errorf("cannot initialize gateway: %v", err)
	}

	session := interactive.InteractiveSessionInit(gateway)
	if generatePaymentIdArgs.issuerToken != "" {
		session.SetIssuerToken(generatePaymentIdArgs.issuerToken)
	}

	if err = session.Login(_registry.Issuer); err != nil {
		return fmt.Errorf("błąd logowania do KSeF: %v", err)
	}

	paymentId, err := session.GeneratePaymentId(seiRefNumbers)
	if err != nil {
		return fmt.Errorf("nie udało się wygenerować identyfikatora płatności: %v", err)
	}

	_registry.PaymentIds = append(_registry.PaymentIds, registry.PaymentId{
		SEIPaymentRefNo: paymentId,
		InvoiceIDS:      seiRefNumbers,
	})

	fmt.Println(paymentId)

	return _registry.Save(generatePaymentIdArgs.path)
}
