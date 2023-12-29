package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
)

func AddMultipartFile(
	multipartWriter *multipart.Writer,
	fieldname string,
	filename string,
	contentType string,
	content io.Reader,
) error {
	header := make(textproto.MIMEHeader)
	header.Set(
		"Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldname, filename),
	)
	header.Set("Content-Type", contentType)
	filePart, err := multipartWriter.CreatePart(header)
	if err != nil {
		return fmt.Errorf("unable to create file writer: %v", err)
	}
	if _, err = io.Copy(filePart, content); err != nil {
		return fmt.Errorf("unable to write file bytes to HTTP request: %v", err)
	}

	return nil
}
