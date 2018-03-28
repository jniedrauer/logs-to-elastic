package parsers

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type EncodedChunk struct {
	Payload []byte
	Records uint32
}

// Iterate through an array in chunks, performing function f
func Chunk(chunkSize int, length int, f func(int, int)) {
	for i := 0; i < length; i += chunkSize {
		end := i + chunkSize
		if end > length {
			end = length
		}

		f(i, end)
	}
}

// Return an encoded chunk of logs
func GetEncodedChunk(chunk []interface{}, delim []byte) ([]byte, error) {
	var enc []byte

	for _, v := range chunk {
		j, err := json.Marshal(v)
		if err != nil {
			log.Error("failed to encode: %v", v)
			continue
		}

		if len(enc) > 0 {
			enc = append(enc, delim...)
		}
		enc = append(enc, j...)
	}

	return enc, nil
}
