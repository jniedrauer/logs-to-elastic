/*
Functions for working with environment variables
*/
package env

import (
	"os"
)

func GetEnvOrDefault(env string, def string) string {
	val, set := os.LookupEnv(env)
	if !set {
		val = def
	}
	return val
}
