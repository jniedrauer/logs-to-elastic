// Equivalent utilities to Python's itertools
package chunk

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
