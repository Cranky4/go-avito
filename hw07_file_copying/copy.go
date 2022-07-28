package main

import (
	"bufio"
	"errors"
	"fmt"
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

func (b *ProgressBar) Done() {
	b.Current = b.Total
}

func (b ProgressBar) Percent() int {
	return int(float64(b.Current) / float64(b.Total) * 100)
}

func (b *ProgressBar) Draw() {
	if b.drawedPercent != b.Percent() {
		b.line = strings.Repeat(b.LineSymbol, b.Percent())
		b.drawedPercent = b.Percent()
		// TODO: terminal clear
		fmt.Printf("[%d%%]%s\n", b.Percent(), b.line)
	}
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

	progressTotal := int(fileInfo.Size())
	if limit > 0 {
		progressTotal = int(math.Min(float64(progressTotal), float64(limit)))
	}

	progressBar := ProgressBar{Total: progressTotal, Current: 0, LineSymbol: "|"}

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

	buffer := make([]byte, bufferSize)

	for {
		count, err := fromFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if limit > 0 {
			count = int(math.Min(float64(count), float64(limit-totalCopied)))
		}

		if count == 0 {
			break
		}

		if writer.Available() < count {
			writer.Flush()
		}

		writer.Write(buffer[0:count])
		progressBar.Advance(count)
		progressBar.Draw()

		if err == io.EOF {
			break
		}

		totalCopied += int64(count)
	}

	// TODO: add terminal clearing
	writer.Flush()
	progressBar.Draw()

	return nil
}
