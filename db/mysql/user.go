package mysql

import (
	"fmt"
	"services/config"
	"services/types"
	"upper.io/db.v3"
)

var userTable string

func init() {
	userTable = config.Config.Database.UserTable
	if userTable == "" {
		panic("user table configuration not set")
	}
}

func (s MySQL) GetUserID(user string) (types.UserCredentials, error) {
	var creds types.UserCredentials

	col := s.DB.Collection(userTable)
	res := col.Find(db.Cond{"username": user})

	cnt, _ := res.Count()
	if cnt == 0 {
		return creds, fmt.Errorf("user %s not found", user)
	}

	defer func() { _ = res.Close() }()

	ok := res.Next(&creds)
	if !ok {
		return creds, fmt.Errorf("user %s not found", user)
	}
	return creds, nil
}
