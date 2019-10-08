package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func (h RestHandlers) fileUpload() error {

	_ = h.r.ParseMultipartForm(32 << 20)

	file, handler, err := h.r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() { _ = file.Close() }()

	fileName, err := h.getParameter("lfs-file")
	if err != nil {
		return err
	}

	fileType, err := h.getParameter("lfs-fileType")
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"fileName": fileName,
		"fileType": fileType,
	}).Debug("Uploading file")

	switch fileType {
	case SURVEY_FILE:
		surveyUpload()
	case GEOG_FILE:
	}

	_, _ = fmt.Fprintf(h.w, "%v", handler.Header)

	f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer func() { _ = f.Close() }()

	_, _ = io.Copy(f, file)

	return nil
}

func surveyUpload() {

}

func (h RestHandlers) getParameter(parameter string) (string, error) {
	keys, ok := h.r.URL.Query()[parameter]

	if !ok || len(keys[0]) < 1 {
		h.log.WithFields(log.Fields{
			"parameter": parameter,
		}).Error("URL parameter missing")
		return "", fmt.Errorf("URL parameter, %s, is missing", parameter)
	}

	return keys[0], nil
}
