package commands

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"ksef/api"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type statusCommand struct {
	Command
}
type statusArgsType struct {
	path   string
	output string
	xml    bool
}

var StatusCommand *statusCommand
var statusArgs statusArgsType

func init() {
	StatusCommand = &statusCommand{
		Command: Command{
			Name:        "status",
			FlagSet:     flag.NewFlagSet("status", flag.ExitOnError),
			Description: "wysyła sprawdza status przesyłki i pobiera dokument UPO",
			Run:         statusRun,
			Args:        statusArgs,
		},
	}

	StatusCommand.FlagSet.StringVar(&statusArgs.path, "p", "", "ścieżka do pliku statusu")
	StatusCommand.FlagSet.StringVar(&statusArgs.output, "o", "", "ścieżka do zapisu UPO (domyślnie katalog pliku statusu + {nrRef}.pdf)")
	StatusCommand.FlagSet.BoolVar(&statusArgs.xml, "xml", false, "zapis UPO jako plik XML")

	registerCommand(&StatusCommand.Command)
}

func statusRun(c *Command) error {
	if statusArgs.path == "" {
		StatusCommand.FlagSet.Usage()
		return nil
	}

	var statusInfo api.StatusInfoFile

	extension := strings.ToLower(filepath.Ext(statusArgs.path))

	if extension != ".yaml" && extension != ".json" {
		return errors.New("unsupported file format. expecting one of .yaml, .json")
	}

	fileContents, err := os.ReadFile(statusArgs.path)

	if err != nil {
		return fmt.Errorf("unable to read status file contents: %v", err)
	}

	if extension == ".yaml" {
		if err = yaml.Unmarshal(fileContents, &statusInfo); err != nil {
			return fmt.Errorf("unable to parse status from yaml: %v", err)
		}
	} else if extension == ".json" {
		if err = json.Unmarshal(fileContents, &statusInfo); err != nil {
			return fmt.Errorf("unable to parse status from json: %v", err)
		}
	}

	if statusInfo.Environment == "" || statusInfo.ReferenceNo == "" {
		return fmt.Errorf("file deserialized correctly, but either environment or referenceNo are empty: %+v", statusInfo)
	}

	gateway, err := api.API_Init(statusInfo.Environment)
	if err != nil {
		return fmt.Errorf("cannot initialize gateway: %v", err)
	}

	var outputFormat = api.UPOFormatPDF
	if statusArgs.xml {
		outputFormat = api.UPOFormatXML
	}

	if statusArgs.output == "" {
		statusArgs.output = filepath.Join(filepath.Dir(statusArgs.path), fmt.Sprintf("%s.%s", statusInfo.ReferenceNo, outputFormat))
	}

	if err = gateway.DownloadUPO(statusInfo.ReferenceNo, outputFormat, statusArgs.output); err != nil {
		return fmt.Errorf("unable to download UPO: %v", err)
	}

	return nil
}
