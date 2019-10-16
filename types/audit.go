package types

import "time"

type Audit struct {
	ReferenceDate time.Time `db:"reference_date"`
	FileName      string    `db:"file_name"`
	NumVarFile    int       `db:"num_var_file"`
	NumVarLoaded  int       `db:"num_var_loaded"`
	NumObFile     int       `db:"num_ob_file"`
	NumObLoaded   int       `db:"num_ob_loaded"`
}
