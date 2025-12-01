package typst

import (
	"bytes"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func (tp *typstPrinter) printTemplate(
	tmpl *template.Template,
	templateName string,
	templateVars any,
	output string,
) (err error) {
	tmpDir, err := os.MkdirTemp(tp.cfg.Workdir, "")
	if err != nil {
		return err
	}
	tmpFile, err := os.CreateTemp(tmpDir, "*.typ")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	if err = tmpl.ExecuteTemplate(tmpFile, templateName, templateVars); err != nil {
		return err
	}

	if tp.cfg.Debug {
		if copyTypErr := utils.CopyFile(tmpFile.Name(), strings.Replace(output, ".pdf", ".typ", 1)); copyTypErr != nil {
			logging.PDFRendererLogger.Error("błąd kopiowania pliku .typ", "err", copyTypErr)
		}
	}

	// now call typst:
	binary := "typst"

	cmd := exec.Command(
		binary,
		"compile", tmpFile.Name(),
	)

	var stdErrBuffer bytes.Buffer

	if tp.cfg.Debug {
		cmd.Stderr = &stdErrBuffer
	}

	defer os.RemoveAll(tmpDir)

	err = cmd.Run()

	if err != nil {
		if tp.cfg.Debug && stdErrBuffer.Len() > 0 {
			if writeErr := os.WriteFile(strings.Replace(output, ".pdf", "-error.txt", 1), stdErrBuffer.Bytes(), 0644); writeErr != nil {
				logging.PDFRendererLogger.Error("błąd zapisywania wyjścia błędu", "err", writeErr)
			}
		}

		return err
	}

	if copyPDFErr := utils.CopyFile(strings.Replace(tmpFile.Name(), ".typ", ".pdf", 1), output); copyPDFErr != nil {
		logging.PDFRendererLogger.Error("błąd kopiowania pliku PDF", "err", copyPDFErr)
	}

	return err

}
