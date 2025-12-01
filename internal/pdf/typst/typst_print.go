package typst

import (
	"ksef/internal/config/pdf/typst"
	"ksef/internal/pdf/printer"
	"ksef/internal/registry/types"
)

type typstPrinter struct {
	cfg *typst.TypstPrinterConfig
}

func NewTypstPrinter(cfg *typst.TypstPrinterConfig) printer.PDFPrinter {
	return &typstPrinter{
		cfg: cfg,
	}
}

func (tp *typstPrinter) PrintUPO(contentBase64 string, output string) error {
	upo, err := printer.ParseUPO(contentBase64)

	if err != nil {
		return err
	}

	return tp.printTemplate(
		tp.cfg.Templates.UPO.Template,
		"upo",
		upo,
		output,
	)
}

func (tp *typstPrinter) Print(contentBase64 string, meta types.Invoice, output string) error {
	return nil
}
