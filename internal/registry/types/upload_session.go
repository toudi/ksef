package types

type InvoiceUploadResult struct {
	Filename string   `yaml:"filename"`
	Checksum string   `yaml:"checksum"`
	SeiRefNo string   `yaml:"seiRefNo"`
	Failed   bool     `yaml:"failed,omitempty"` // whether the invoice was processed successfuly
	Errors   []string `yaml:"errors,omitempty"`
}

type UploadSessionStatus struct {
	Processed bool                   `yaml:"processed"`
	Invoices  []*InvoiceUploadResult `yaml:"invoices"`
	UPO       []string               `yaml:"upo,omitempty"`
}
