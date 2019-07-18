package services

import (
	"fmt"
	"hash"
	"log"
	"pds-go/lfs/db"
	_ "pds-go/lfs/db"
	"pds-go/lfs/db/model"
)

func GetUsers() []model.USER {
	var rr []model.USER
	var conn = db.DB

	e := conn.Find(&rr)
	if e.Error != nil {
		log.Fatal(fmt.Errorf("cannot read users table %v", e))
	}

	return rr
}

type Password struct {
	Digest     func() hash.Hash
	SaltSize   int
	KeyLen     int
	Iterations int
}
