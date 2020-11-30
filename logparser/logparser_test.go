package logparser_test

import (
	"apps/logparser"
	"github.com/parmaanu/goutils/filesystem"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	mfs := filesystem.NewMockFileSystem()
	filesystem.SetFileSystem(mfs)

	mfs.SetFileData("file1.txt", []string{
		"and gradually we made this country more just and more equal",
		"and that office should take care of all the people and be the custodian of this office",
		"one constitutional office elected by the people is the presidency",
		"regardless of religion, ideals and social status will be protected",
	})

	mfs.SetFileData("file2.txt", []string{
		"a system of representative goverenment",
		"i am in philedalphia, where our consitution was written and signed",
		"the right to participate in the political process",
		"tonight I want to talk as plainly as I can.",
	})

	fname := "test.log"
	fdata := []string{"line1",
		"tradelog, orderId=[123]",
		"line3",
		"line4, tradelog, orderId=[456]",
		"line5, insert order sent, price:1234, quantity",
		"line6, insert orde sent, price:1234, quantity",
		"line7, insert order reply, price:1234, quantity",
	}
	mfs.SetFileData(fname, fdata)

	fdata1 := []string{}
	fdata2 := []string{}
	for i := 0; i < 100; i++ {
		line1 := "ORDER qty:" + strconv.Itoa(rand.Intn(100))
		line2 := "EXEC qty:" + strconv.Itoa(rand.Intn(100))
		fdata1 = append(fdata1, line1)
		fdata2 = append(fdata2, line2)
		if i%10 == 0 {
			fdata1 = append(fdata1, "this line should be counted")
		}
	}
	mfs.SetFileData("file_large1.txt", fdata1)
	mfs.SetFileData("file_large2.txt", fdata2)
}

func TestLogparser(t *testing.T) {
	matchedLinesOrder := 0
	matchedLinesExec := 0

	lpr := logparser.NewLogParser()

	lpr.AddFileSources("file_large1.txt", "file_large2.txt")

	// Test single pattern matching
	lpr.AddConfig(logparser.Config{
		Patterns: []string{"ORDER"},
		OnEachLineFunc: func(config *logparser.OnEachLineConfig) {
			assert.Equal(t, "ORDER", config.Pat, "pattern not matching in OnEachLineFunc function")
			matchedLinesOrder++
		},
	})

	lpr.AddConfig(logparser.Config{
		Patterns: []string{"EXEC"},
		OnEachLineFunc: func(config *logparser.OnEachLineConfig) {
			assert.Equal(t, "EXEC", config.Pat, "pattern not matching in OnEachLineFunc function")
			matchedLinesExec++
		},
	})

	lpr.Run()

	assert.Equal(t, 100, matchedLinesOrder, "Number of orders don't match")
	assert.Equal(t, 100, matchedLinesExec, "Number of execs don't match")
}

func TestMultiSourceLineReader(t *testing.T) {
	mslr := logparser.NewMultiSourceLineReader()

	file1, _ := logparser.NewFileLineReader("file1.txt")
	file2, _ := logparser.NewFileLineReader("file2.txt")

	mslr.AddSources(file1, file2)

	expectedLines := []string{
		"a system of representative goverenment",
		"and gradually we made this country more just and more equal",
		"and that office should take care of all the people and be the custodian of this office",
		"i am in philedalphia, where our consitution was written and signed",
		"one constitutional office elected by the people is the presidency",
		"regardless of religion, ideals and social status will be protected",
		"the right to participate in the political process",
		"tonight I want to talk as plainly as I can.",
	}

	readLines := []string{}
	for {
		nextLine, err := mslr.NextLine()
		if err == -1 {
			break
		}
		readLines = append(readLines, nextLine)
	}

	assert.Equal(t, expectedLines, readLines)
}
