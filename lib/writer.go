package lib

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ResultWriter provides buffered, thread-safe file output.
type ResultWriter struct {
	mu     sync.Mutex
	writer *bufio.Writer
	file   *os.File
}

// NewResultWriter creates a buffered writer for the given path.
func NewResultWriter(path string) (*ResultWriter, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}
	if dir != "." {
		if err := os.Chmod(dir, 0o700); err != nil {
			return nil, fmt.Errorf("chmod directory: %w", err)
		}
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	if err := f.Chmod(0o600); err != nil {
		return nil, errors.Join(fmt.Errorf("chmod file: %w", err), f.Close())
	}
	return &ResultWriter{
		file:   f,
		writer: bufio.NewWriter(f),
	}, nil
}

// WriteLine appends a line to the buffer.
func (rw *ResultWriter) WriteLine(line string) error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	_, err := fmt.Fprintln(rw.writer, line)
	return err
}

// Close flushes the buffer and closes the file.
func (rw *ResultWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return errors.Join(rw.writer.Flush(), rw.file.Close())
}
