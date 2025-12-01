package latex

import (
	"ksef/internal/pdf/printer"
)

func (lp *LatexPrinter) PrintUPO(contentBase64 string, output string) error {
	upo, err := printer.ParseUPO(contentBase64)

	if err != nil {
		return err
	}

	return lp.printTemplate(
		lp.cfg.Templates.UPO.Template,
		"upo",
		upo,
		output,
	)
}
