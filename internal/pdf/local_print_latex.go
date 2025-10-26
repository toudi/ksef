package pdf

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"ksef/internal/config/pdf/latex"
	"ksef/internal/logging"
	"ksef/internal/registry"
	"ksef/internal/sei/generators/fa_3_1"
	"ksef/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type LatexPrinter struct {
	cfg *latex.LatexPrinterConfig
}

type QrCodesInfo struct {
	Invoice     string
	Certificate string
}

type LatexTemplateVars struct {
	Invoice  *Invoice
	Registry registry.Invoice
	Qrcodes  QrCodesInfo
	Header   latex.HeaderFooterSettings
	Footer   latex.HeaderFooterSettings
}

func (lp *LatexPrinter) Print(contentBase64 string, meta registry.Invoice, output string) error {
	var invoiceXMLBuffer bytes.Buffer
	invoiceXMLBytes, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return err
	}
	if _, err = invoiceXMLBuffer.Write(invoiceXMLBytes); err != nil {
		return err
	}

	var i = Invoice{arrayElements: fa_3_1.FA_3_1ArrayElements}
	if err = xml.Unmarshal(invoiceXMLBytes, &i); err != nil {
		return err
	}

	var templateData = LatexTemplateVars{
		Invoice:  &i,
		Registry: meta,
		Qrcodes: QrCodesInfo{
			Invoice: meta.SEIQRCode,
		},
		Header: lp.cfg.Templates.Invoice.Header,
		Footer: lp.cfg.Templates.Invoice.Footer,
	}

	tmpDir, err := os.MkdirTemp(lp.cfg.Workdir, "")
	if err != nil {
		return err
	}
	tmpFile, err := os.CreateTemp(tmpDir, "*.tex")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	if err = lp.cfg.Templates.Invoice.Template.ExecuteTemplate(tmpFile, "invoice", templateData); err != nil {
		return err
	}

	// now call pdflatex:
	cmd := exec.Command(
		"docker",
		"run", "--rm", "--name", "latex", "-v", tmpDir+":/workdir",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"texlive/texlive",
		"pdflatex", filepath.Join("/workdir", filepath.Base(tmpFile.Name())),
	)

	var stdErrBuffer bytes.Buffer

	if lp.cfg.Debug {
		cmd.Stderr = &stdErrBuffer
	}

	defer os.RemoveAll(tmpDir)

	err = cmd.Run()

	if err != nil {
		if lp.cfg.Debug && stdErrBuffer.Len() > 0 {
			if writeErr := os.WriteFile(strings.Replace(output, ".pdf", "-error.txt", 1), stdErrBuffer.Bytes(), 0644); writeErr != nil {
				logging.PDFRendererLogger.Error("błąd zapisywania wyjścia błędu", "err", writeErr)
			}
		}

		return err
	}

	if lp.cfg.Debug {
		if copyTexErr := utils.CopyFile(tmpFile.Name(), strings.Replace(output, ".pdf", ".tex", 1)); copyTexErr != nil {
			logging.PDFRendererLogger.Error("błąd kopiowania pliku .tex", "err", copyTexErr)
		}
	}

	if copyPDFErr := utils.CopyFile(strings.Replace(tmpFile.Name(), ".tex", ".pdf", 1), output); copyPDFErr != nil {
		logging.PDFRendererLogger.Error("błąd kopiowania pliku PDF", "err", copyPDFErr)
	}

	return err
}
