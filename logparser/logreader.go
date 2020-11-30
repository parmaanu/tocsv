package logparser

import (
	"bufio"
	"github.com/parmaanu/goutils/filesystem"
	"io"
)

// LineReader defines a interfaces for reading from a source line by line
type LineReader interface {
	NextLine() (string, error)
	GetCurrentLineNumber() int
	Finished() bool
}

// FileLineReader implements a LineReader interface for Files
type FileLineReader struct {
	UnderlyingFile filesystem.File
	Reader         *bufio.Reader
	LineNo         int
	EOFReached     bool
}

// TODO: Support gzip files
// Will do as part of issue #10
func (flr *FileLineReader) open(filepath string) error {
	file, err := filesystem.Open(filepath)
	if err != nil {
		return err
	}
	flr.UnderlyingFile = file
	flr.Reader = bufio.NewReader(flr.UnderlyingFile)
	flr.EOFReached = false
	return nil
}

// Close the underlying file
func (flr *FileLineReader) Close() {
	flr.UnderlyingFile.Close()
}

// NextLine returns the nextline along with error
func (flr *FileLineReader) NextLine() (string, error) {
	line, err := flr.Reader.ReadString('\n')
	if err != nil {
		// Couldn't find the new line delimiter, maybe EOF
		if err == io.EOF && len(line) > 0 {
			// This is valid case with some data in line, which needs to be processed
			flr.LineNo++
			return line, err
		}
		// All these other cases are to be treated as this file is no longer valid to be
		// read
		flr.EOFReached = true
		return "", err
	}
	flr.LineNo++
	return line[:len(line)-1], err
}

// GetCurrentLineNumber returns numbers of lines read till now
func (flr *FileLineReader) GetCurrentLineNumber() int {
	return flr.LineNo
}

// Finished returns true if EOF Reached
func (flr *FileLineReader) Finished() bool {
	return flr.EOFReached
}

// NewFileLineReader returns an object FileLineReader
// If file opening fails, it returns nil, err
func NewFileLineReader(filename string) (*FileLineReader, error) {
	flr := &FileLineReader{}
	err := flr.open(filename)
	if err != nil {
		return nil, err
	}
	return flr, nil
}

// MultiSourceLineReader contains multiple sources of type Interface LineReader
// and returns lines in lexicographically sorted manner from among the sources
type MultiSourceLineReader struct {
	Sources     []LineReader
	CurrentLine []string
}

// AddSources adds to the list of sources
func (mslr *MultiSourceLineReader) AddSources(linereaders ...LineReader) {
	for index := range linereaders {
		mslr.Sources = append(mslr.Sources, linereaders[index])
		currentLine, _ := linereaders[index].NextLine()
		mslr.CurrentLine = append(mslr.CurrentLine, currentLine)
	}
}

// NextLine returns NextLine and a status code
// status code = 0, Success
// status code = 1, Failure
func (mslr *MultiSourceLineReader) NextLine() (nextLine string, err int) {
	minIndex := -1
	for i := 0; i < len(mslr.Sources); i++ {
		if mslr.Sources[i].Finished() {
			continue
		}
		if (minIndex == -1) || (mslr.CurrentLine[i] < nextLine) {
			nextLine = mslr.CurrentLine[i]
			minIndex = i
		}
	}
	if minIndex == -1 {
		return "", -1
	}
	mslr.CurrentLine[minIndex], _ = mslr.Sources[minIndex].NextLine()
	return nextLine, 0
}

// NewMultiSourceLineReader create a new MultiLineSourceReader
func NewMultiSourceLineReader() *MultiSourceLineReader {
	return &MultiSourceLineReader{}
}
