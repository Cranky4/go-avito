package main

import (
	"bufio"
	"errors"
	"io"
	"math"
	"os"
	"strings"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidParameter      = errors.New("some parameters is invalid")
)

type ProgressBar struct {
	Current, Total, drawedPercent int
	LineSymbol, line              string
}

func (b *ProgressBar) Advance(number int) {
	b.Current += number
}

func (b ProgressBar) Percent() int {
	return int(float64(b.Current) / float64(b.Total) * 100)
}

func (b *ProgressBar) Draw() string {
	if b.drawedPercent != b.Percent() {
		b.line = strings.Repeat(b.LineSymbol, b.Percent())
		b.drawedPercent = b.Percent()
	}

	return b.line
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromPath)

	if limit < 0 || offset < 0 {
		return ErrInvalidParameter
	}

	if err != nil || fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	progressBar := ProgressBar{Total: int(fileInfo.Size()), Current: 0, LineSymbol: "|"}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer fromFile.Close()
	fromFile.Seek(offset, 0)

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	writer := bufio.NewWriter(toFile)
	var totalCopied int64
	var bufferSize uint64 = 128
	if limit > 0 {
		bufferSize = uint64(math.Min(float64(bufferSize), float64(limit)))
	}

	for {
		if limit > 0 {
			bufferSize = uint64(math.Min(float64(bufferSize), float64(limit-totalCopied)))
		}

		buffer := make([]byte, bufferSize)
		count, err := fromFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		if count == 0 {
			break
		}

		if writer.Available() < count {
			writer.Flush()
		}

		writer.Write(buffer[0:count])
		// fmt.Printf("%s\n\n\n\n", string(buffer))
		progressBar.Advance(count)
		// fmt.Printf("[%d%%]%s\n", progressBar.Percent(), progressBar.Draw())

		if err == io.EOF {
			break
		}

		totalCopied += int64(count)
	}

	writer.Flush()

	return nil
}
