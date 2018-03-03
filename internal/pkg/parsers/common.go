package parsers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"

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
func GetEncodedChunk(start int, end int, delim []byte, f func(int, int) []interface{}) []byte {
	var enc []byte

	for _, v := range f(start, end) {
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

	return enc
}

// Get number of newlines in a Reader
func LineCount(r io.Reader) int {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
			log.Fatalf(err.Error())
		}
	}
}

// Get a number of lines starting at a byte offset
func GetLines(start int64, lines int, data io.ReadSeeker) ([][]byte, int64) {
	var output [][]byte

	if _, err := data.Seek(start, 0); err != nil {
		log.Fatalf(err.Error())
	}

	scanner := bufio.NewScanner(data)

	var offset int64 = 0
	for i := 0; i < lines; i++ {
		if scanner.Scan() {
			bytes := scanner.Bytes()
			offset += int64(len(bytes))
			output := append(output, bytes)
		} else if err := scanner.Err(); err != nil {
			log.Fatalf(err.Error())
		} else {
			break
		}
	}

	return output, offset
}
