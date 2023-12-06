package inputprocessors

import (
	"fmt"
	"ksef/common"
	"strings"

	"github.com/xuri/excelize/v2"
)

type XLSXDecoder struct {
	config xlsxConfig
}

func XLSXDecoder_Init(config xlsxConfig) *XLSXDecoder {
	return &XLSXDecoder{config: config}
}

func (x *XLSXDecoder) Process(sourceFile string, parser *common.Parser) error {
	// let's check if the file exists
	var err error
	var file *excelize.File
	file, err = excelize.OpenFile(sourceFile)

	if err != nil {
		return err
	}
	defer file.Close()

	if file.SheetCount == 0 {
		return fmt.Errorf("There are no sheets in this workbook")
	}

	// let's select the first sheet ..
	var sheetName = file.GetSheetName(0)
	// unless there's a particular one that we should inspect
	if x.config.SheetName != "" {
		if _, err = file.GetSheetIndex(x.config.SheetName); err != nil {
			return fmt.Errorf("Specified sheet does not seem to exist within the file: %v", err)
		}
		sheetName = x.config.SheetName
	}

	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("error fetching sheet rows: %v", err)
	}

	for _, row := range rows {
		if len(row) >= 2 && strings.ToLower(row[1]) == common.SectionInvoice {
			if err = parser.ClearInvoiceData(); err != nil {
				return err
			}
		}

		err = parser.ProcessLine(row)
		if err != nil {
			return err
		}
	}

	// notify the parser that we've finished processing the file so there is
	// definetely one unprocessed invoice within the file
	return parser.InvoiceReady()

}
