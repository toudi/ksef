package batch

import (
	"bytes"
	"fmt"
	"ksef/internal/logging"
	registryPkg "ksef/internal/registry"
	"ksef/internal/sei/api/client"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

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
}

var batchInitResponse batchInitResponseType
var finishResponsePayload finishResponseType

const (
	endpointBatchInit   = "/api/batch/Init"
	endpointBatchFinish = "/api/batch/Finish"
)

func (b *BatchSession) UploadInvoices(sourcePath string) error {
	var log *slog.Logger = logging.BatchLogger
	signedMetadataFile, err := locateBatchMetadataFile(sourcePath)

	if err != nil {
		return fmt.Errorf("error locating metadata file: %v", err)
	}

	session := client.NewHTTPSession(b.apiClient.Environment.Host)

	// step 1 - send initAuthRequest
	_, err = session.SendFile(client.SendFileParams{
		Method:       "POST",
		Endpoint:     endpointBatchInit,
		FileName:     signedMetadataFile,
		JSONResponse: &batchInitResponse,
		Logger:       logging.BatchHTTPLogger,
	})
	if err != nil {
		return fmt.Errorf("unable to send file: %v", err)
	}

	// // step 2 - upload encrypted archive
	log.Debug(
		"BatchSession::UploadInvoices",
		"url",
		batchInitResponse.PackageSignature.PackagePartSignatureList[0].Url,
		"method", "PUT",
	)

	encryptedArchive, err := os.Open(filepath.Join(sourcePath, "metadata.zip.aes"))
	if err != nil {
		return fmt.Errorf("could not open encrypted archive for sending: %v", err)
	}
	defer encryptedArchive.Close()
	stat, _ := encryptedArchive.Stat()

	batchPutRequest, err := http.NewRequest(
		"PUT",
		batchInitResponse.PackageSignature.PackagePartSignatureList[0].Url,
		encryptedArchive,
	)
	if err != nil {
		return fmt.Errorf("could not prepare PUT request: %v", err)
	}
	batchPutRequest.Header.Set("Content-Type", "application/octet-stream")
	batchPutRequest.Header.Set("accept", "application/json")
	batchPutRequest.ContentLength = stat.Size()

	for _, header := range batchInitResponse.PackageSignature.PackagePartSignatureList[0].HeaderEntryList {
		logging.BatchHTTPLogger.Debug(
			"BatchSession::UploadInvoices",
			"set header",
			"header",
			header.Key,
			"value",
			header.Value,
		)

		batchPutRequest.Header.Set(header.Key, header.Value)
	}
	batchResponse, err := http.DefaultClient.Do(batchPutRequest)
	if err != nil {
		return fmt.Errorf("could not send encrypted archive: %v", err)
	}
	defer batchResponse.Body.Close()
	// batchResponseBody, err := io.ReadAll(batchResponse.Body)
	if err != nil {
		return fmt.Errorf("could not read batch upload response: %v", err)
	}
	// fmt.Printf("result of PUT request: %d\n", batchResponse.StatusCode)
	// fmt.Printf("%s\n", string(batchResponseBody))
	if batchResponse.StatusCode/100 != 2 {
		return fmt.Errorf("unexpected response code from PUT request.")
	}

	// step 3 - call finish upload
	finishResponsePayload.ReferenceNumber = batchInitResponse.ReferenceNumber

	finishResponse, err := session.JSONRequest(
		client.JSONRequestParams{
			Method:   "POST",
			Endpoint: endpointBatchFinish,
			Payload:  finishResponsePayload,
			Response: nil,
			Logger:   logging.BatchHTTPLogger,
		},
	)

	if err != nil {
		return fmt.Errorf("could not call finish request: %v", err)
	}

	// fmt.Printf("result of batch/Finish call: %d\n", finishResponse.StatusCode)

	if finishResponse.StatusCode/100 != 2 {
		return fmt.Errorf("bad response from finishUpload: %d", finishResponse.StatusCode)
	}

	// step 4 - persist status for fetching UPO.
	registry := registryPkg.NewRegistry()
	registry.Environment = b.apiClient.EnvironmentAlias
	registry.SessionID = batchInitResponse.ReferenceNumber
	return registry.Save(path.Join(sourcePath, "registry.yaml"))
}

func locateBatchMetadataFile(sourcePath string) (string, error) {
	files, err := os.ReadDir(sourcePath)
	if err != nil {
		return "", fmt.Errorf("unable to read dir: %v", err)
	}

	for _, fileInfo := range files {
		var fileName = fileInfo.Name()

		if filepath.Ext(fileName) == ".xml" {
			// check if this file contains the Signature part
			xmlContents, err := os.ReadFile(filepath.Join(sourcePath, fileName))
			if err != nil {
				return "", fmt.Errorf("unable to read invoice file: %v", err)
			}
			if bytes.Contains(xmlContents, []byte(":Signature>")) {
				return filepath.Join(sourcePath, fileName), nil
			}
		}
	}

	return "", fmt.Errorf("unable to find batch metadata file")
}
