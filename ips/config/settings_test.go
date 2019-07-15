package config

import (
	"testing"
)

func TestServer(t *testing.T) {
	server := Config.Database.Server
	if server != "localhost" {
		t.Errorf("server = %s; want localhost", server)
	} else {
		t.Logf("Server %s\n", server)
	}

	user := Config.Database.User
	if user != "ips" {
		t.Errorf("user = %s; want ips", user)
	} else {
		t.Logf("user %s\n", user)
	}

	password := Config.Database.Password
	if password != "ips" {
		t.Errorf("password = %s; want ips", password)
	} else {
		t.Logf("password %s\n", password)
	}

	databaseName := Config.Database.Database
	if databaseName != "ips" {
		t.Errorf("database name = %s; want ips", databaseName)
	} else {
		t.Logf("database name %s\n", databaseName)
	}

	debug := Config.Debug
	if debug != true {
		t.Errorf("debug = %t; want localhost", debug)
	} else {
		t.Logf("debug %t\n", debug)
	}
}
