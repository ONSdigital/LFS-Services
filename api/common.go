package api

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"services/db"
	"services/types"
	"services/util"
	"time"
)

const (
	Error = "ERROR"
	OK    = "OK"
)

func SaveStreamToTempFile(w http.ResponseWriter, r *http.Request) (string, error) {

	file, _, err := r.FormFile("lfsFile")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error getting formfile")
		return "", err
	}

	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()

	_ = r.ParseMultipartForm(64 << 20)

	fileName := r.Form.Get("fileName")
	if fileName == "" {
		log.Error().Msg("fileName not set")
		return "", fmt.Errorf("fileName not set")
	}

	log.Debug().
		Str("fileName", fileName).
		Msg("Uploading file")

	startTime := time.Now()

	tmpfile, err := ioutil.TempFile("", fileName)
	if err != nil {
		return "", fmt.Errorf("cannot create temporary file: %s ", err)
	}

	n, err := io.Copy(tmpfile, file)

	log.Debug().
		Str("fileName", fileName).
		Int64("bytesRead", n).
		Str("elapsedTime", util.FmtDuration(startTime)).
		Msg("File uploaded")

	_ = tmpfile.Close()
	return tmpfile.Name(), nil
}

func FindNIBatch(monthNo, yearNo int) (types.NIBatchItem, error) {

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot connect to database")
		return types.NIBatchItem{}, fmt.Errorf("cannot connect to database: %s", err)
	}

	info, err := database.FindNIBatchInfo(monthNo, yearNo)
	if err != nil {
		log.Error().Err(err)
		return types.NIBatchItem{},
			fmt.Errorf("cannot upload the survey file as the batch for month %d and year %d does not exist", monthNo, yearNo)
	}

	return info, nil
}

func FindGBBatch(weekNo, yearNo int) (types.GBBatchItem, error) {

	database, err := db.GetDefaultPersistenceImpl()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Cannot connect to database")
		return types.GBBatchItem{}, fmt.Errorf("cannot connect to database: %s", err)
	}

	info, err := database.FindGBBatchInfo(weekNo, yearNo)
	if err != nil {
		log.Error().Err(err)
		return types.GBBatchItem{},
			fmt.Errorf("cannot upload the survey file as the batch for week %d and year %d does not exist",
				weekNo, yearNo)
	}

	return info, nil
}
