package commands

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type statusCommand struct {
	Command
}
type statusArgsType struct {
	path string
}

var StatusCommand *statusCommand
var statusArgs statusArgsType

type statusResponseType struct {
	ProcessingCode  int    `json:"processingCode"`
	Upo             string `json:"upo"`
	ReferenceNumber string `json:"ReferenceNumber"`
}

var statusResponse statusResponseType

func init() {
	StatusCommand = &statusCommand{
		Command: Command{
			Name:        "status",
			FlagSet:     flag.NewFlagSet("status", flag.ExitOnError),
			Description: "wysyła sprawdza status przesyłki i pobiera dokument UPO",
			Run:         statusRun,
			Args:        statusArgs,
		},
	}

	StatusCommand.FlagSet.StringVar(&statusArgs.path, "p", "", "ścieżka do pliku statusu")

	registerCommand(&StatusCommand.Command)
}

func statusRun(c *Command) error {
	if statusArgs.path == "" {
		StatusCommand.FlagSet.Usage()
		return nil
	}

	var workdir = filepath.Dir(statusArgs.path)

	statusUrlBytes, err := ioutil.ReadFile(statusArgs.path)
	if err != nil {
		return fmt.Errorf("could not read %s: %v", statusArgs.path, err)
	}

	statusReq, err := http.NewRequest("GET", string(statusUrlBytes), nil)
	if err != nil {
		return fmt.Errorf("could not prepare GET request: %v", err)
	}
	statusHTTPResponse, err := http.DefaultClient.Do(statusReq)
	if err != nil || statusHTTPResponse.StatusCode/100 != 2 {
		return fmt.Errorf("could not execute HTTP request: %v", err)
	}
	defer statusHTTPResponse.Body.Close()
	err = json.NewDecoder(statusHTTPResponse.Body).Decode(&statusResponse)
	if err != nil {
		return fmt.Errorf("cannot decode HTTP response to JSON struct: %v", err)
	}

	fmt.Printf("GET %s returned status %d\n", statusReq.URL, statusHTTPResponse.StatusCode)

	upoXml, _ := base64.StdEncoding.DecodeString(statusResponse.Upo)

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", statusResponse.ReferenceNumber+".xml"))
	h.Set("Content-Type", "text/xml")
	part, _ := writer.CreatePart(h)

	io.Copy(part, bytes.NewReader(upoXml))
	writer.Close()

	htmlUpoReq, err := http.NewRequest("POST", fmt.Sprintf("https://%s/web/api/session/get-upo-html-view", statusReq.Host), body)
	htmlUpoReq.Header.Set("Content-Type", writer.FormDataContentType())

	htmlUpoResponse, err := http.DefaultClient.Do(htmlUpoReq)
	if err != nil || htmlUpoResponse.StatusCode/100 != 2 {
		content, _ := ioutil.ReadAll(htmlUpoResponse.Body)
		fmt.Printf("%s\n", string(content))
		htmlUpoResponse.Body.Close()
		return fmt.Errorf("cannot fetch HTML UPO: %d / %v", htmlUpoResponse.StatusCode, err)
	}
	defer htmlUpoResponse.Body.Close()

	var upoHtmlFilename = filepath.Join(workdir, statusResponse.ReferenceNumber+"-upo.html")
	var upoPDFFilename = filepath.Join(workdir, statusResponse.ReferenceNumber+"-upo.pdf")

	htmlUpo, err := os.Create(upoHtmlFilename)
	if err != nil {
		return fmt.Errorf("cannot create HTML upo file: %v", err)
	}
	defer htmlUpo.Close()

	// base64 decode the upo:
	upoBase64Content, err := ioutil.ReadAll(htmlUpoResponse.Body)
	if err != nil {
		return fmt.Errorf("cannot read base64-encoded upo: %v", err)
	}
	upoHtmlContent, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(string(upoBase64Content), `"`, ""))
	if err != nil {
		return fmt.Errorf("cannot decode base64 upo: %v", err)
	}
	htmlUpo.WriteString(string(upoHtmlContent))

	// check if we can covert it to pdf
	var canConvertToPDF bool = false
	var converter string = "wkhtmltopdf"

	if runtime.GOOS == "windows" {
		converter = "wkhtmltopdf.exe"
		_, err = os.Stat(converter)
		canConvertToPDF = !os.IsNotExist(err)
	} else {
		_cmd, err := exec.LookPath("wkhtmltopdf")
		canConvertToPDF = (_cmd != "" && err == nil)
	}

	if canConvertToPDF {
		cmd := exec.Command(converter, upoHtmlFilename, upoPDFFilename)
		return cmd.Run()
	}
	return nil
}
