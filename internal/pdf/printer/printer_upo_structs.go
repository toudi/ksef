package printer

import (
	"encoding/xml"
	"time"
)

type Session struct {
	RefNo string `xml:"NumerReferencyjnySesji"`
}

type LogicalStruct struct {
	Name     string `xml:"NazwaStrukturyLogicznej"`
	FormCode string `xml:"KodFormularza"`
}

type AuthContext struct {
	Type  string
	Value string
}

type Auth struct {
	Context  AuthContext `xml:"IdKontekstu"`
	Checksum string      `xml:"SkrotDokumentuUwierzytelniajacego"`
}

type Confirmation struct {
	Page          int `xml:"Strona"`
	TotalPages    int `xml:"LiczbaStron"`
	DocumentsFrom int `xml:"ZakresDokumentowOd"`
	DocumentsTo   int `xml:"ZakresDokumentowDo"`
	Documents     int `xml:"CalkowitaLiczbaDokumentow"`
}

type UPOItem struct {
	KSeFRefNo          string    `xml:"NumerKSeFDokumentu"`
	RefNo              string    `xml:"NumerFaktury"`
	SellerNIP          string    `xml:"NipSprzedawcy"`
	IssueDate          string    `xml:"DataWystawieniaFaktury"`
	SendDate           time.Time `xml:"DataPrzeslaniaDokumentu"`
	NumberAssignedDate time.Time `xml:"DataNadaniaNumeruKSeF"`
	Checksum           string    `xml:"SkrotDokumentu"`
}

type UPO struct {
	XMLName xml.Name `xml:"Potwierdzenie"`
	Session
	Confirmation Confirmation `xml:"OpisPotwierdzenia"`
	LogicalStruct
	Auth     Auth      `xml:"Uwierzytelnienie"`
	Invoices []UPOItem `xml:"Dokument"`
}
