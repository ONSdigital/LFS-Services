package config_test

import (
	"fmt"
	conf "services/config"
	"testing"
)

func TestConfig(t *testing.T) {
	server := conf.Config.Database.Server
	if server != "localhost" {
		t.Errorf("server = %s; want localhost", server)
	} else {
		t.Logf("Server %s\n", server)
	}

	user := conf.Config.Database.User
	if user != "lfs" {
		t.Errorf("user = %s; want lfs", user)
	} else {
		t.Logf("user %s\n", user)
	}

	password := conf.Config.Database.Password
	if password != "lfs" {
		t.Errorf("password = %s; want lfs", password)
	} else {
		t.Logf("password %s\n", password)
	}

	databaseName := conf.Config.Database.Database
	if databaseName != "LFS" {
		t.Errorf("database name = %s; want LFS", databaseName)
	} else {
		t.Logf("database name %s\n", databaseName)
	}

	maxPoolsize := conf.Config.Database.ConnectionPool.MaxPoolSize
	if maxPoolsize != 10 {
		t.Errorf("maxPoolsize = %d; want 10", maxPoolsize)
	} else {
		t.Logf("maxPoolsize %d\n", maxPoolsize)
	}

	columnsTable := conf.Config.Database.SurveyTable
	if columnsTable != "columns" {
		t.Errorf("columnsTable = %s, want columns", columnsTable)
	} else {
		t.Logf("columnsTable %s\n", columnsTable)
	}

	address := conf.Config.Service.ListenAddress
	expected := "127.0.0.1:8000"
	if address != expected {
		t.Errorf("address = %s, want %s", address, expected)
	} else {
		t.Logf("port %s\n", address)
	}

	readTimeout := conf.Config.Service.ReadTimeout
	expectedTimeout := "60s"
	if readTimeout != expectedTimeout {
		t.Errorf("readTimeout = %s, want %s", readTimeout, expectedTimeout)
	} else {
		t.Logf("readTimeout %s\n", readTimeout)
	}

	writeTimeout := conf.Config.Service.WriteTimeout
	expectedTimeout = "60s"
	if readTimeout != expectedTimeout {
		t.Errorf("writeTimeout = %s, want %s", writeTimeout, expectedTimeout)
	} else {
		t.Logf("writeTimeout %s\n", writeTimeout)
	}

	ren := conf.Config.Rename.Survey

	for _, v := range ren {
		fmt.Printf("Rename from: %s, to: %s\n", v.From, v.To)
	}

	drop := conf.Config.DropColumns.Survey

	for _, v := range drop.ColumnNames {
		fmt.Printf("Drop Column: %s\n", v)
	}
}
