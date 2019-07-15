package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml"
)

var Config configuration

func init() {
	configFile, err := ioutil.ReadFile(fileName())

	if err != nil {
		log.Fatal(fmt.Errorf("cannot read configuration %v", err))
	}

	Config = configuration{}

	err = toml.Unmarshal(configFile, &Config)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot unmarshall configuration file %v", err))
	}

	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.Info("Configuration loaded")

}

func fileName() string {
	env := os.Getenv("ENV")

	if len(env) == 0 {
		env = "development"
	}

	filename := []string{"config.", env, ".toml"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
