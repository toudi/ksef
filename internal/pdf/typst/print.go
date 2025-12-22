package typst

import (
	"bytes"
	"io/fs"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (tp *typstPrinter) prepareWorkdir() error {
	if err := os.MkdirAll(tp.cfg.Workdir, 0775); err != nil {
		return err
	}
	_, err := os.Stat(filepath.Join(tp.cfg.Workdir, "upo"))
	if os.IsNotExist(err) {
		// copy everything from templates-dir to workdir
		if err = filepath.WalkDir(tp.cfg.Templates, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			destLocalPath := strings.TrimPrefix(path, tp.cfg.Templates)
			if destLocalPath == "" {
				return nil
			}

			if d.IsDir() {
				return os.MkdirAll(filepath.Join(tp.cfg.Workdir, destLocalPath), 0775)
			}

			return utils.CopyFile(
				path,
				filepath.Join(tp.cfg.Workdir, destLocalPath),
			)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (tp *typstPrinter) print(
	template string,
	output string,
) (err error) {
	cmd := exec.Command(
		"typst",
		"compile",
		"--root",
		tp.cfg.Workdir,
		filepath.Join(tp.cfg.Workdir, template),
		output,
	)

	logging.PDFRendererLogger.Debug("executing", "cmd", cmd.String())

	var stdErrBuffer bytes.Buffer

	if tp.cfg.Debug {
		cmd.Stderr = &stdErrBuffer
	}

	err = cmd.Run()
	if err != nil {
		if tp.cfg.Debug && stdErrBuffer.Len() > 0 {
			if writeErr := os.WriteFile(strings.Replace(output, ".pdf", "-error.txt", 1), stdErrBuffer.Bytes(), 0644); writeErr != nil {
				logging.PDFRendererLogger.Error("błąd zapisywania wyjścia błędu", "err", writeErr)
			}
		}

		return err
	}

	return nil
}
