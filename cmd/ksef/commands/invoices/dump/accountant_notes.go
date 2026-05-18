package dump

import (
	"bytes"
	"fmt"
	"ksef/internal/logging"
	"ksef/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// AnnotationsFile is the YAML structure written for the typst accountant-notes template.
type AnnotationsFile struct {
	Metadata    AnnotationsMetadata `yaml:"metadata"`
	Annotations []AnnotationEntry   `yaml:"annotations"`
}

// AnnotationsMetadata holds the report metadata.
type AnnotationsMetadata struct {
	ReportDate string `yaml:"report-date"`
	Generator  string `yaml:"generator"`
}

// AnnotationEntry is a single annotation row in the YAML file.
type AnnotationEntry struct {
	Seller   string `yaml:"seller,omitempty"`
	Invoice  string `yaml:"invoice,omitempty"`
	Item     string `yaml:"item,omitempty"`
	ItemName string `yaml:"item-name,omitempty"`
	Notes    string `yaml:"notes,omitempty"`
}

// renderAccountantNotes creates a YAML file from the collected annotations,
// renders the accountant-notes.typ Typst template into a PDF, and saves it
// to the specified output path.
func renderAccountantNotes(invoiceAnnotations []InvoiceAnnotations, outputPath string) error {
	// Build the YAML data structure
	yamlData := AnnotationsFile{
		Metadata: AnnotationsMetadata{
			ReportDate: time.Now().Local().Format("2006-01-02"),
			Generator:  dumpAnnotationCfg.Generator,
		},
	}

	for _, inv := range invoiceAnnotations {
		for _, itemRule := range inv.ItemRules {
			var notes string
			if itemRule.Annotation != nil {
				notes = itemRule.Annotation.String()
			}

			yamlData.Annotations = append(yamlData.Annotations, AnnotationEntry{
				Seller:   inv.SellerName,
				Invoice:  inv.RefNo,
				Item:     itemRule.Item.OrdNo,
				ItemName: itemRule.Item.Name,
				Notes:    notes,
			})
		}
	}

	// Create a temporary directory for the typst work
	workDir, err := os.MkdirTemp("", "accountant-notes-*")
	if err != nil {
		return fmt.Errorf("błąd tworzenia katalogu tymczasowego: %w", err)
	}
	defer os.RemoveAll(workDir)

	// Write the YAML file to the temp directory
	yamlPath := filepath.Join(workDir, "annotations.yaml")
	if err := utils.SaveYAML(yamlData, yamlPath); err != nil {
		return fmt.Errorf("błąd zapisywania pliku YAML: %w", err)
	}

	// Copy the YAML file into the template directory (so typst can find it via
	// the path configured in the config)
	if err := utils.CopyFile(yamlPath, dumpAnnotationCfg.YamlPath); err != nil {
		return fmt.Errorf("błąd kopiowania pliku YAML do szablonu: %w", err)
	}

	// Compile the Typst template
	typstOutputPath := filepath.Join(workDir, "accountant-notes.pdf")
	cmdExec := exec.Command(
		"typst",
		"compile",
		dumpAnnotationCfg.TemplatePath,
		typstOutputPath,
	)

	var stdErrBuffer bytes.Buffer
	cmdExec.Stderr = &stdErrBuffer

	if err := cmdExec.Run(); err != nil {
		logging.PDFRendererLogger.Error("błąd kompilacji typst", "err", err, "stderr", stdErrBuffer.String())
		return fmt.Errorf("błąd kompilacji typst: %w", err)
	}

	// Copy the generated PDF to the output path
	if err := utils.CopyFile(typstOutputPath, outputPath); err != nil {
		return fmt.Errorf("błąd kopiowania pliku PDF: %w", err)
	}

	logging.InvoicesDBLogger.Info("wygenerowano accountant-notes.pdf", "plik", outputPath)
	return nil
}
