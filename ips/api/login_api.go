package api

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	p "github.com/wuriyanto48/go-pbkdf2"
	"net/http"
	"strings"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	password := vars["password"]

	db := loginDatabase()
	log.Info(fmt.Printf("Checking User: %s, password: %s\n", vars["user"], vars["password"] ))
	getUsers(user, password, db)

	w.WriteHeader(http.StatusOK)

}

func loginDatabase() *sql.DB {
	connectionString := "ips:ips@tcp(127.0.0.1:3306)/ips"
	var db, err = sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(fmt.Errorf("login failed %v", err))
	}

	return db
}

func getUsers(user string, password string, db *sql.DB) {

	var (
		userName     string
		userPassword string
	)

	rows, err := db.Query("select user_name, password from user where USER_NAME = ?", user)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot read users table %v", err))
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Warn(fmt.Errorf("received error on rows.Close() %v", err))
		}
	}()
	
	for rows.Next() {
		err := rows.Scan(&userName, &userPassword)

		if err != nil {
			log.Error(fmt.Errorf("users table empty? %v", err))
		}

		s := strings.Split(userPassword, "$")
		salt := s[1]
		pw := s[2]
		log.Println("password: ", pw)

		pass := p.NewPassword(sha256.New, len(salt), len(pw), 150000)
		hashed := p.HashResult{CipherText: pw, Salt: salt}

		isValid := pass.VerifyPassword(password, hashed.CipherText, hashed.Salt)

		log.Println("Password is valid: ", isValid)
	}
}
