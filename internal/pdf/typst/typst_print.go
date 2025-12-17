package typst

import (
	"ksef/internal/config/pdf/typst"
	"ksef/internal/pdf/printer"
	"ksef/internal/registry/types"
	"path/filepath"
)

type typstPrinter struct {
	cfg *typst.TypstPrinterConfig
}

func NewTypstPrinter(cfg *typst.TypstPrinterConfig) printer.PDFPrinter {
	return &typstPrinter{
		cfg: cfg,
	}
}

func (tp *typstPrinter) PrintUPO(srcFile string, output string) error {
	var templateVars = map[string]string{
		"upoXML": filepath.Base(srcFile),
	}
	return tp.printTemplate(
		tp.cfg.Templates.UPO.Template,
		"upo",
		templateVars,
		output,
		srcFile,
	)
}

func (tp *typstPrinter) Print(contentBase64 string, meta types.Invoice, output string) error {
	return nil
}
