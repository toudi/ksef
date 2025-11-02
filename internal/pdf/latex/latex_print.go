package latex

import (
	"bytes"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// abstract function that prints using the LaTeX engine
func (lp *LatexPrinter) printTemplate(
	tmpl *template.Template,
	templateName string,
	templateVars any,
	output string,
) (err error) {
	tmpDir, err := os.MkdirTemp(lp.cfg.Workdir, "")
	if err != nil {
		return err
	}
	tmpFile, err := os.CreateTemp(tmpDir, "*.tex")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	if err = tmpl.ExecuteTemplate(tmpFile, templateName, templateVars); err != nil {
		return err
	}

	if lp.cfg.Debug {
		if copyTexErr := utils.CopyFile(tmpFile.Name(), strings.Replace(output, ".pdf", ".tex", 1)); copyTexErr != nil {
			logging.PDFRendererLogger.Error("błąd kopiowania pliku .tex", "err", copyTexErr)
		}
	}

	// now call pdflatex:
	runner := "docker"
	if lp.cfg.Podman {
		runner = "podman"
	}
	cmd := exec.Command(
		runner,
		"run", "--rm", "--name", "latex", "-v", tmpDir+":/workdir",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"texlive/texlive",
		"pdflatex", filepath.Join("/workdir", filepath.Base(tmpFile.Name())),
	)

	var stdErrBuffer bytes.Buffer

	if lp.cfg.Debug {
		cmd.Stdout = &stdErrBuffer
	}

	defer os.RemoveAll(tmpDir)

	err = cmd.Run()

	if err != nil {
		stdErrBuffer.Cap()
		fmt.Printf("err: %+v\n; stderr: %+v\n", err, stdErrBuffer.Cap())
		if lp.cfg.Debug && stdErrBuffer.Len() > 0 {
			if writeErr := os.WriteFile(strings.Replace(output, ".pdf", "-error.txt", 1), stdErrBuffer.Bytes(), 0644); writeErr != nil {
				logging.PDFRendererLogger.Error("błąd zapisywania wyjścia błędu", "err", writeErr)
			}
		}

		return err
	}

	if copyPDFErr := utils.CopyFile(strings.Replace(tmpFile.Name(), ".tex", ".pdf", 1), output); copyPDFErr != nil {
		logging.PDFRendererLogger.Error("błąd kopiowania pliku PDF", "err", copyPDFErr)
	}

	return err

}
