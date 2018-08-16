// Load configuration from environment variables
package conf

import (
	"log"
	"os"
	"strconv"
)

const defaultDelimiter = "\n"
const defaultTimeout = 10

// TODO: This whole thing is garbage

// Stores configuration state
type Config struct {
	ChunkSize       int
	Delimiter       []byte
	IndexName       string
	Logstash        string
	LogstashTimeout int
}

// Populate configuration state
func NewConfig() *Config {
	d := []byte(GetEnvOrDefault("DELIMITER", defaultDelimiter))
	i := GetEnvOrFatal("INDEXNAME")
	l := GetEnvOrFatal("LOGSTASH")

	cs, err := strToInt(GetEnvOrFatal("CHUNK_SIZE"))
	if err != nil {
		log.Fatalf("not a number: %v", err)
	}
	lt, err := strToInt(GetEnvOrDefault("LOGSTASH_TIMEOUT", strconv.Itoa(defaultTimeout)))
	if err != nil {
		log.Fatalf("not a number: %v", err)
	}

	return &Config{Delimiter: d, IndexName: i, Logstash: l, ChunkSize: cs, LogstashTimeout: lt}
}

func strToInt(str string) (result int, err error) {
	i64, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return
	}

	return int(i64), nil
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
