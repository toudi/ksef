package subjectsettings

import (
	"errors"
	"ksef/internal/invoicesdb/jpk/constants"
	"slices"
)

type JPKFormMeta struct {
	IRSCode    int               `yaml:"irs-code,omitempty"`
	SystemName string            `yaml:"system-name,omitempty"`
	Subject    *SubjectData      `yaml:"subject,omitempty"`
	Defaults   map[string]string `yaml:"defaults,omitempty"`
}

type SubjectData struct {
	SubjectType string            `yaml:"type"`
	Data        map[string]string `yaml:"data"`
}

type SurplusAction struct {
	CarryOver bool   `yaml:"carry-over"`
	Refund    string `yaml:"refund,omitempty"`
	OffsetTax string `yaml:"offset-tax,omitempty"`
}

type JPKSettings struct {
	FormMeta JPKFormMeta   `yaml:"form,omitempty"`
	Surplus  SurplusAction `yaml:"surplus,omitempty"`
}

func (s JPKSettings) Validate() error {
	if !s.Surplus.CarryOver && s.Surplus.Refund == "" && s.Surplus.OffsetTax == "" {
		return errors.New("nie wybrano trybu rozliczenia nadwyżki VAT")
	}
	if s.Surplus.CarryOver && (s.Surplus.Refund != "" || s.Surplus.OffsetTax != "") {
		return errors.New("opcje przeniesienia na następny okres rozliczeniowy oraz zwrotu na rachunek lub przeksięgowania są wzajemnie rozłączne")
	}
	if !s.Surplus.CarryOver && (s.Surplus.Refund == "" && s.Surplus.OffsetTax == "") || (s.Surplus.Refund != "" && s.Surplus.OffsetTax != "") {
		return errors.New("opcja zwrotu na rachunek oraz przeksięgowania na poczet przyszłych zobowiązań są wzajemnie rozłączne")
	}
	if s.Surplus.Refund != "" {
		if !slices.Contains(constants.VATRefundModes, s.Surplus.Refund) {
			return errors.New("niedozwolona wartość trybu zwrotu nadwyżki")
		}
	}

	return nil
}
