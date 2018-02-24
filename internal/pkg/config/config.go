/*
Load configuration from an SSM parameter

TODO: Documentation
*/
package config

import (
	"fmt"
	"strings"

	"github.com/op/go-logging"
	//"gopkg.in/yaml.v2"
	//"github.com/aws/aws-sdk-go/service/ssm"
)

var log = logging.MustGetLogger("config")

type Configurationer interface {
	LoadConfig()
}

type Configuration struct {
	ssmConfigParam   string
	LogstashEndpoint string `yaml:"logstash"`
	LogLevel         string `yaml:"log_level"`
	LogGroups        string `yaml:"log_groups"`
}

func (f Configuration) LoadConfig() {
	fmt.Println(strings.ToUpper("gopher"))
}

type LogGroup struct {
	LogGroup  string `yaml:"logGroup"`
	IndexName string `yaml:"indexname"`
}
