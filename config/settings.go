package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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

	switch Config.LogFormat {
	case "Text":
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	case "Json":
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat:  "",
			DisableTimestamp: false,
			DataKey:          "",
			FieldMap:         nil,
			CallerPrettyfier: nil,
			PrettyPrint:      false,
		})
	case "Terminal":
		break
	}

	switch Config.LogLevel {
	case "Trace":
		log.SetLevel(log.TraceLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	case "Fatal":
		log.SetLevel(log.FatalLevel)

	}

}

func fileName() string {
	runEnv := os.Getenv("CONFIG")

	if len(runEnv) == 0 {
		runEnv = "development"
	}

	filename := []string{"config.", runEnv, ".toml"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
