// The purpose of this file is
// 1) to show how to make a NewReader from strings package to implement a Readcloser interface.
// 2) to demonstrate the danger of not closing files — accumulating open file handles
// can cause resource exhaustion of the OS and failure to open more files.
// Note that the program tries to open a text file somefile.txt, defined in line 20.
//
// To make the leak apparent when not closing the file descriptors, you can execute with:
// `ulimit -n 256 && go run .`

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	const file = "somefile.txt" // any small readable file
	const limit = 512

	// Example 1: strings.Reader (in-memory)
	// No need to close the reader.
	for i := 0; i < limit; i++ {
		str := fmt.Sprintf("hello#%d", i)
		inmem := &ReaderWithCloser{
			Reader: strings.NewReader(str),
		}
		use(inmem)
	}

	// Example 2: os.File
	// Real file, must be closed.
	for i := 0; i < limit; i++ {
		f, err := os.Open("somefile.txt")
		if err != nil {
			fmt.Printf("\n❌ Failed at %d files: %v\n", i, err)
			break
		}
		file := &ReaderWithCloser{
			Reader: f,
		}
		use(file)
	}
}

type ReaderWithCloser struct {
	Reader io.Reader
}

func (r *ReaderWithCloser) Read(p []byte) (int, error) {
	return r.Reader.Read(p)
}

func (r *ReaderWithCloser) Close() error {
	if c, ok := r.Reader.(io.Closer); ok {
		return c.Close()
	}
	return nil // safe no-op for readers without Close
}

func use(rc io.ReadCloser) {
	buf := make([]byte, 20)

	n, _ := rc.Read(buf)
	fmt.Printf("Read: %q\n", buf[:n])

	// Uncomment the line below to prevent leaking resources.
	// rc.Close()
}
