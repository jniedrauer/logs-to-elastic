package parsers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"

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

// Get number of newlines in a file
func LineCount(fileName string) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	file, err := os.Open(fileName)
	if err != nil {
		return count, err
	}
	defer file.Close()

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			if count <= 0 {
				return 1, errors.New("no newline found")
			} else {
				return count, nil
			}

		case err != nil:
			return count, err
		}
	}
}

// Get a number of lines starting at a byte offset
func GetLines(start int64, lines int, fileName string) ([][]byte, int64, error) {
	var output [][]byte

	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		return output, start, err
	}

	if _, err := file.Seek(start, io.SeekStart); err != nil {
		return output, start, err
	}

	scanner := bufio.NewScanner(file)

	var offset int64 = 0
	for i := 0; i < lines; i++ {
		if scanner.Scan() {
			bytes := scanner.Bytes()
			offset += int64(len(bytes) + 1) // 1 here is for the newline byte
			output = append(output, bytes)
		} else if err := scanner.Err(); err != nil {
			return [][]byte{}, offset, err
		} else {
			break
		}
	}

	return output, offset, nil
}
