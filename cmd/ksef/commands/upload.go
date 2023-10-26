package commands

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type uploadCommand struct {
	Command
}

type uploadArgsType struct {
	testGateway bool
	path        string
}

var UploadCommand *uploadCommand
var uploadArgs = &uploadArgsType{}

const testGatewayURL = "https://ksef-test.mf.gov.pl/api/"
const productionGatewayURL = "https://ksef.mf.gov.pl/api/"

type batchInitResponseType struct {
	ReferenceNumber  string `json:"referenceNumber"`
	PackageSignature struct {
		PackagePartSignatureList []struct {
			HeaderEntryList []struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			} `json:"headerEntryList"`
			Url string `json:"url"`
		} `json:"packagePartSignatureList"`
	} `json:"packageSignature"`
}

type finishResponseType struct {
	ReferenceNumber string `json:"referenceNumber"`
	//Timestamp       string `json:"timestamp"`
}

var batchInitResponse batchInitResponseType
var finishResponse finishResponseType

func init() {
	UploadCommand = &uploadCommand{
		Command: Command{
			Name:        "upload",
			FlagSet:     flag.NewFlagSet("upload", flag.ExitOnError),
			Description: "wysyła podpisany plik KSEF do bramki ministerstwa finansów",
			Run:         uploadRun,
			Args:        uploadArgs,
		},
	}

	UploadCommand.FlagSet.BoolVar(&uploadArgs.testGateway, "t", false, "użyj bramki testowej")
	UploadCommand.FlagSet.StringVar(&uploadArgs.path, "p", "", "ścieżka do katalogu z wygenerowanymi fakturami")

	registerCommand(&UploadCommand.Command)
}

func uploadRun(c *Command) error {
	if uploadArgs.path == "" {
		c.FlagSet.Usage()
		return nil
	}

	var gateway = productionGatewayURL
	if uploadArgs.testGateway {
		gateway = testGatewayURL
	}

	var workdir = filepath.Dir(uploadArgs.path)

	// step 1 - upload metadata

	metadataFile, err := os.Open(uploadArgs.path)
	if err != nil {
		return fmt.Errorf("could not open input file for sending: %v", err)
	}
	defer metadataFile.Close()

	fmt.Printf("step 1 - POST to %v\n", gateway+"batch/Init")
	initRequest, err := http.NewRequest("POST", gateway+"batch/Init", metadataFile)
	if err != nil {
		return fmt.Errorf("error preparing request: %v", err)
	}
	initResponse, err := http.DefaultClient.Do(initRequest)
	if err != nil {
		return fmt.Errorf("error sending file: %v", err)
	}
	defer initResponse.Body.Close()
	fmt.Printf("initResponse status: %d\n", initResponse.StatusCode)
	if initResponse.StatusCode/100 != 2 {
		responseContent, err := io.ReadAll(initResponse.Body)
		if err != nil {
			return fmt.Errorf("error reading response:%v", err)
		}
		return fmt.Errorf("unexpected response code from initResponse:\n%s\n", string(responseContent))
	}

	err = json.NewDecoder(initResponse.Body).Decode(&batchInitResponse)
	if err != nil {
		return fmt.Errorf("error decoding JSON response: %v", err)
	}

	fmt.Printf("batch init response: %+v\n", batchInitResponse)

	// step 2 - upload encrypted archive
	fmt.Printf("step 2 - PUT to %v\n", batchInitResponse.PackageSignature.PackagePartSignatureList[0].Url)
	encryptedArchive, err := os.Open(filepath.Join(workdir, "metadata.zip.aes"))
	if err != nil {
		return fmt.Errorf("could not open encrypted archive for sending: %v", err)
	}
	defer encryptedArchive.Close()
	stat, _ := encryptedArchive.Stat()

	batchPutRequest, err := http.NewRequest("PUT", batchInitResponse.PackageSignature.PackagePartSignatureList[0].Url, encryptedArchive)
	if err != nil {
		return fmt.Errorf("could not prepare PUT request: %v", err)
	}
	batchPutRequest.Header.Set("Content-Type", "application/octet-stream")
	batchPutRequest.Header.Set("accept", "application/json")
	batchPutRequest.ContentLength = stat.Size()

	for _, header := range batchInitResponse.PackageSignature.PackagePartSignatureList[0].HeaderEntryList {
		fmt.Printf("add header %v with a value of %v\n", header.Key, header.Value)
		batchPutRequest.Header.Set(header.Key, header.Value)
	}
	batchResponse, err := http.DefaultClient.Do(batchPutRequest)
	if err != nil {
		return fmt.Errorf("could not send encrypted archive: %v", err)
	}
	defer batchResponse.Body.Close()
	batchResponseBody, err := io.ReadAll(batchResponse.Body)
	if err != nil {
		return fmt.Errorf("could not read batch upload response: %v", err)
	}
	fmt.Printf("result of PUT request: %d\n", batchResponse.StatusCode)
	fmt.Printf("%s\n", string(batchResponseBody))
	if batchResponse.StatusCode/100 != 2 {
		return fmt.Errorf("unexpected response code from PUT request.")
	}

	// step 3 - call finish upload

	finishResponse.ReferenceNumber = batchInitResponse.ReferenceNumber
	var finishResponseBuffer bytes.Buffer
	if err = json.NewEncoder(&finishResponseBuffer).Encode(finishResponse); err != nil {
		return fmt.Errorf("cannot encode finishResponse to JSON")
	}

	fmt.Printf("step 3 - POST to %v\n", gateway+"batch/Finish")
	finishUpload, err := http.NewRequest("POST", gateway+"batch/Finish", &finishResponseBuffer)
	finishUpload.Header.Set("Content-Type", "application/json")

	if err != nil {
		return fmt.Errorf("could not prepare finish request: %v", err)
	}
	finishResponse, err := http.DefaultClient.Do(finishUpload)
	if err != nil {
		return fmt.Errorf("could not finish upload: %v", err)
	}
	defer finishResponse.Body.Close()
	responseBody, err := io.ReadAll(finishResponse.Body)
	if err != nil {
		return fmt.Errorf("could not read response from finishUpload: %v", err)
	}

	fmt.Printf("result of batch/Finish call: %d\n", finishResponse.StatusCode)
	fmt.Printf("%s\n", string(responseBody))

	if finishResponse.StatusCode/100 != 2 {
		return fmt.Errorf("bad response from finishUpload: %v", responseBody)
	}

	// step 4 - persist the url for fetching UPO.
	return os.WriteFile(filepath.Join(workdir, "metadata.ref"), []byte(fmt.Sprintf("%scommon/Status/%s", gateway, batchInitResponse.ReferenceNumber)), 0644)
}
