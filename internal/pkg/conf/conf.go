package conf

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	Logstash  string
	IndexName string
	ChunkSize int
}

func Init() *Conf {
	l := os.Getenv("LOGSTASH")
	i := os.Getenv("INDEXNAME")

	c64, err := strconv.ParseInt(os.Getenv("CHUNK_SIZE"), 10, 0)
	if err != nil {
		log.Fatalf(err.Error())
	}
	cs := int(c64)

	return &Conf{Logstash: l, IndexName: i, ChunkSize: cs}
}
