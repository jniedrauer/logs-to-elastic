// Load configuration from environment variables
package conf

import (
	"log"
	"os"
	"strconv"
)

var defaultRegion string = "us-east-1"

// Stores configuration state
type Conf struct {
	Logstash  string
	IndexName string
	ChunkSize int
	AwsRegion string
}

// Populate configuration state
func Init() *Conf {
	l := GetEnvOrFatal("LOGSTASH")
	i := GetEnvOrFatal("INDEXNAME")

	c64, err := strconv.ParseInt(GetEnvOrFatal("CHUNK_SIZE"), 10, 0)
	if err != nil {
		log.Fatalf("not a number: %v", err)
	}
	cs := int(c64)

	return &Conf{Logstash: l, IndexName: i, ChunkSize: cs}
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
