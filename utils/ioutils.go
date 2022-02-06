package utils

import (
	"errors"
	"github.com/jfrog/build-info-go/entities"
	"io"
	"os"
	"strings"
	"sync"
)

var ErrShortWrite = errors.New("short write")

type asyncMultiWriter struct {
	writers []io.Writer
}

// AsyncMultiWriter creates a writer that duplicates its writes to all the
// provided writers asynchronous
func AsyncMultiWriter(writers ...io.Writer) io.Writer {
	w := make([]io.Writer, len(writers))
	copy(w, writers)
	return &asyncMultiWriter{w}
}

// Writes data asynchronously to each writer and waits for all of them to complete.
// In case of an error, the writing will not complete.
func (t *asyncMultiWriter) Write(p []byte) (int, error) {
	var wg sync.WaitGroup
	wg.Add(len(t.writers))
	errChannel := make(chan error)
	finished := make(chan bool, 1)
	for _, w := range t.writers {
		go writeData(p, w, &wg, errChannel)
	}
	go func() {
		wg.Wait()
		close(finished)
	}()
	// This select will block until one of the two channels returns a value.
	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

func writeData(p []byte, w io.Writer, wg *sync.WaitGroup, errChan chan error) {
	n, err := w.Write(p)
	if err != nil {
		errChan <- err
	}
	if n != len(p) {
		errChan <- ErrShortWrite
	}
	wg.Done()
}

func UnixToWinPathSeparator(filePath string) string {
	return strings.Replace(filePath, "/", "\\\\", -1)
}

func WinToUnixPathSeparator(filePath string) string {
	return strings.Replace(filePath, "\\", "/", -1)
}

func GetFileDetails(filePath string, includeChecksums bool) (*FileDetails, error) {
	var err error
	details := new(FileDetails)
	if includeChecksums {
		details.Checksum, err = calcChecksumDetails(filePath)
		if err != nil {
			return details, err
		}
	} else {
		details.Checksum = entities.Checksum{}
	}

	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	details.Size = fileInfo.Size()
	return details, nil
}
