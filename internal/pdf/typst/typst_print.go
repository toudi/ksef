package typst

import (
	"ksef/internal/config/pdf/typst"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/logging"
	"ksef/internal/pdf/printer"
	"ksef/internal/utils"
	"path/filepath"
)

type typstPrinter struct {
	cfg *typst.TypstPrinterConfig
}

func Printer(cfg *typst.TypstPrinterConfig) printer.PDFPrinter {
	return &typstPrinter{
		cfg: cfg,
	}
}

func (tp *typstPrinter) PrintUPO(srcFile string, output string) error {
	if err := tp.prepareWorkdir(); err != nil {
		return err
	}

	templatePath := filepath.Dir(
		filepath.Join(
			tp.cfg.Workdir, tp.cfg.UPO.Template,
		),
	)

	destPath := filepath.Join(templatePath, "upo.xml")
	logging.PDFRendererLogger.Debug(
		"copy", "src", srcFile, "dest", destPath,
	)
	if err := utils.CopyFile(
		srcFile,
		destPath,
	); err != nil {
		return err
	}

	return tp.print(
		tp.cfg.UPO.Template,
		output,
	)
}

func (tp *typstPrinter) PrintInvoice(
	srcFile string,
	output string,
	meta *monthlyregistry.InvoicePrintingMeta,
) error {
	if err := tp.prepareWorkdir(); err != nil {
		return err
	}

	templatePath := filepath.Dir(
		filepath.Join(
			tp.cfg.Workdir, tp.cfg.Invoice.Template,
		),
	)

	destPath := filepath.Join(templatePath, "invoice.xml")
	metaYAML := filepath.Join(templatePath, "meta.yaml")

	meta.Page.Header.Left = tp.cfg.Invoice.Header.Left
	meta.Page.Header.Center = tp.cfg.Invoice.Header.Center
	meta.Page.Header.Right = tp.cfg.Invoice.Header.Right

	if err := utils.SaveYAML(meta, metaYAML); err != nil {
		return err
	}

	if err := utils.CopyFile(
		srcFile,
		destPath,
	); err != nil {
		return err
	}

	return tp.print(
		tp.cfg.Invoice.Template,
		output,
	)
}
