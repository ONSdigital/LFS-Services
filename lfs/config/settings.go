package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
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
		log.Fatal(fmt.Errorf("cannot read configuration %+v", err))
	}

	Config = configuration{}

	err = toml.Unmarshal(configFile, &Config)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot unmarshall configuration file %+v", err))
	}

	// Parse environment variables
	if err := env.Parse(&Config); err != nil {
		log.Fatal(fmt.Errorf("cannot parse environment variables %+v", err))
	}

	if Config.Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	log.Info("Configuration loaded")

}

func fileName() string {
	runEnv := os.Getenv("ENV")

	if len(runEnv) == 0 {
		runEnv = "development"
	}

	filename := []string{"config.", runEnv, ".toml"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
