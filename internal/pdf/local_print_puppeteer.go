package pdf

import (
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/registry/types"
	"os/exec"
	"strings"
)

type PuppeteerReferencePrinter struct {
	nodeBin         string
	browserBin      string
	templatePath    string
	renderingScript string
}

func (p *PuppeteerReferencePrinter) Print(
	contentBase64 string,
	invoiceMeta types.Invoice,
	output string,
) error {
	renderedFilePath, err := preparePrerenderedTemplate(
		p.templatePath,
		&invoiceMeta,
		contentBase64,
	)
	if err != nil {
		return fmt.Errorf("unable to prepare prerendered file: %v", err)
	}
	command := p.nodeBin
	commandArgs := []string{
		p.renderingScript,
		p.browserBin,
		renderedFilePath,
		output,
	}

	logging.PDFRendererLogger.Debug("local PDF render", "command", command, "args", commandArgs)

	cmd := exec.Command(command, commandArgs...)
	b := new(strings.Builder)
	cmd.Stderr = b

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error from stderr: %v\n", b)
		return fmt.Errorf("error during conversion to PDF: %v", err)
	}

	return nil
}
