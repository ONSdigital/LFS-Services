package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tkanos/gonfig"
)

var (
	Config     Configuration
	DBServer   string
	DBUser     string
	DBPassword string
	DBName     string
	Debug      bool
)

func init() {

	Config = Configuration{}
	err := gonfig.GetConf(getFileName(), &Config)

	if err != nil {
		log.Fatal(fmt.Errorf("cannot read configurstion %v", err))
	}

	DBServer = Config.Database.Server
	DBUser = Config.Database.User
	DBPassword = Config.Database.Password
	DBName = Config.Database.Database
	Debug = Config.Debug

	if Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.Info("Configuration loaded")

}

func getFileName() string {
	env := os.Getenv("ENV")

	if len(env) == 0 {
		env = "development"
	}

	filename := []string{"config.", env, ".json"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
