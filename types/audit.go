package types

import "time"

type AuditStatus int

const (
	NotStarted   AuditStatus = 0
	FileUploaded AuditStatus = 1
	FileReloaded AuditStatus = 2
	UploadFailed AuditStatus = 3
)

type Audit struct {
	Id            int         `db:"id" json:"id"`
	FileName      string      `db:"file_name" json:"fileName"`
	FileSource    FileSource  `db:"file_source" json:"fileSource"`
	Week          int         `db:"week" json:"week"`
	Month         int         `db:"month" json:"month"`
	Year          int         `db:"year" json:"year"`
	ReferenceDate time.Time   `db:"reference_date" json:"referenceDate"`
	NumVarFile    int         `db:"num_var_file" json:"numVarFile"`
	NumVarLoaded  int         `db:"num_var_loaded" json:"numVarLoaded"`
	NumObFile     int         `db:"num_ob_file" json:"numObFile"`
	NumObLoaded   int         `db:"num_ob_loaded" json:"numObLoaded"`
	Status        AuditStatus `db:"status" json:"status"`
	Message       string      `db:"message" json:"message"`
}

type ErrorResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"errorMessage"`
}

type OkayResponse struct {
	Status string `json:"status"`
}
