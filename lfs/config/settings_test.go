package config_test

import (
	conf "pds-go/lfs/config"
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
	if databaseName != "lfs" {
		t.Errorf("database name = %s; want lfs", databaseName)
	} else {
		t.Logf("database name %s\n", databaseName)
	}

	debug := conf.Config.Debug
	if debug != true {
		t.Errorf("debug = %t; want localhost", debug)
	} else {
		t.Logf("debug %t\n", debug)
	}
}
