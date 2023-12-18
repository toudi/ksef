package inputprocessors

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"ksef/internal/sei/constants"
	"ksef/internal/sei/parser"
	"os"
	"strconv"
	"strings"
)

type CSVFormat struct {
	delimiter              string
	encodingConversionFile string
	encodingConversion     map[byte]string
}

func CSVDecoder_Init(config csvConfig) (InputProcessor, error) {
	decoder := &CSVFormat{
		delimiter:              config.Delimiter,
		encodingConversionFile: config.EncodingConversionFile,
	}

	if decoder.encodingConversionFile != "" {
		decoder.prepareEncodingConversionTable()
	}

	return decoder, nil
}

func (c *CSVFormat) Process(sourceFile string, parser *parser.Parser) error {
	// let's check if the file exists
	var err error
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return err
	}
	csvFile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	// we cannot use the "regular" csv reader because it assumes that each line has
	// the same number of fields which is not applicable here.

	scanner := bufio.NewScanner(csvFile)

	var lineNo = 1

	for scanner.Scan() {
		line := scanner.Text()

		csvReader := csv.NewReader(strings.NewReader(line))
		if c.delimiter != "" {
			csvReader.Comma = rune(c.delimiter[0])
		}

		fields, err := csvReader.Read()
		if err != nil {
			return fmt.Errorf("error during reading CSV at line %d: %v", lineNo, err)
		}

		if c.encodingConversion != nil {
			for i, field := range fields {
				fields[i] = c.convertEncoding(field)
			}
		}

		// check if this is a new invoice in batch mode
		if len(fields) >= 2 && strings.ToLower(fields[1]) == constants.SectionInvoice {
			if err = parser.ClearInvoiceData(); err != nil {
				return err
			}
		}

		err = parser.ProcessLine(fields)
		if err != nil {
			return err
		}

		lineNo += 1
	}

	// notify the parser that we've finished processing the file so there is
	// definetely one unprocessed invoice within the file
	return parser.InvoiceReady()
}

func (c *CSVFormat) convertEncoding(data string) string {
	if c.encodingConversion == nil {
		c.prepareEncodingConversionTable()
	}

	inputBytes := []byte(data)
	outputBytes := []byte{}
	var dstByte []byte

	for _, srcByte := range inputBytes {
		dstByte = []byte{srcByte}
		if dstChar, exists := c.encodingConversion[srcByte]; exists {
			dstByte = bytes.Replace(dstByte, []byte{srcByte}, []byte(dstChar), 1)
		}
		outputBytes = append(outputBytes, dstByte...)
	}

	return string(outputBytes)

}

func (c *CSVFormat) prepareEncodingConversionTable() {
	fileBytes, err := os.ReadFile(c.encodingConversionFile)

	if err != nil {
		// log.Errorf("error opening encoding conversion file")
		return
	}

	c.encodingConversion = make(map[byte]string)

	for _, line := range strings.Split(string(fileBytes), LineBreak) {
		mapping := strings.Split(line, ":")

		if len(mapping) == 2 {
			srcByteHex := strings.Trim(mapping[0], " ")
			dstChar := strings.Trim(mapping[1], " \r")

			srcByte, err := strconv.ParseUint(srcByteHex, 0, 8)
			if err == nil {
				c.encodingConversion[byte(srcByte)] = dstChar
			}
		}
	}
}
