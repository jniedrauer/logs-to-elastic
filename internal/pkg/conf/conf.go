// Load configuration from environment variables
package conf

import (
	"log"
	"os"
	"strconv"
)

var defaultRegion string = "us-east-1"
var defaultDelimiter string = "\n"

// Stores configuration state
type Config struct {
	AwsRegion string
	ChunkSize int
	Delimiter []byte
	IndexName string
	Logstash  string
}

// Populate configuration state
func NewConfig() *Config {
	a := GetEnvOrDefault("AWS_REGION", defaultRegion)
	d := []byte(GetEnvOrDefault("DELIMITER", defaultDelimiter))
	i := GetEnvOrFatal("INDEXNAME")
	l := GetEnvOrFatal("LOGSTASH")

	c64, err := strconv.ParseInt(GetEnvOrFatal("CHUNK_SIZE"), 10, 0)
	if err != nil {
		log.Fatalf("not a number: %v", err)
	}
	cs := int(c64)

	return &Config{AwsRegion: a, Delimiter: d, IndexName: i, Logstash: l, ChunkSize: cs}
}

// Get environment variable or return default if unset
func GetEnvOrDefault(env string, def string) string {
	val, set := os.LookupEnv(env)
	if !set {
		val = def
	}
	return val
}

// Get environment variable or raise a fatal error
func GetEnvOrFatal(env string) string {
	val, set := os.LookupEnv(env)
	if !set {
		log.Fatalf("Environment variable not set: %s", env)
	}

	return val
}
