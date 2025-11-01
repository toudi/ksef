package pdf

import (
	"fmt"
	"ksef/internal/registry"
	"os"
	"path"
	"strings"
)

func preparePrerenderedTemplate(
	templateFile string,
	invoiceMeta *registry.Invoice,
	contentBase64 string,
) (string, error) {
	replacer := strings.NewReplacer(
		"__qrcode_url__", invoiceMeta.QRCodes.Invoice,
		"__invoice_sei_ref_no__", invoiceMeta.KSeFReferenceNumber,
		"__invoice_base64__", contentBase64,
	)

	templateDir := path.Dir(templateFile)
	templateContent, err := os.ReadFile(templateFile)
	if err != nil {
		return "", fmt.Errorf("unable to read template file: %v", err)
	}

	renderFilePath := path.Join(templateDir, "render.html")
	destFile, err := os.Create(renderFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to create pre-rendered file: %v", err)
	}

	_, err = replacer.WriteString(destFile, string(templateContent))
	destFile.Close()
	if err != nil {
		return "", fmt.Errorf("unable to write to destFile: %v", err)
	}

	return renderFilePath, nil
}
