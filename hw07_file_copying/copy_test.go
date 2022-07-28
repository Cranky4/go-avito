package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func compareFiles(from, to string) bool {
	fromInfo, _ := os.Stat(from)
	inputBuf := make([]byte, fromInfo.Size())
	input, _ := os.Open(from)
	io.ReadFull(input, inputBuf)

	toInfo, _ := os.Stat(from)
	outputBuf := make([]byte, toInfo.Size())
	output, _ := os.Open(to)
	io.ReadFull(output, outputBuf)

	result := bytes.Compare(inputBuf, outputBuf)

	input.Close()
	output.Close()

	return result == 0
}

var (
	fromFile = "./testdata/input.txt"
	toFile   = "./testdata/output.txt"
)

func TestCopy(t *testing.T) {
	fromFileInfo, _ := os.Stat(fromFile)

	t.Run("big offset - err", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 1000000, 0)

		require.NotNil(t, err)
		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "offset exceeds file size", err)
		_, fileInfoErr := os.Stat(toFile)
		require.True(t, os.IsNotExist(fileInfoErr))
	})

	t.Run("no limit, no offset", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 0, 0)

		require.Nil(t, err)
		toFileInfo, fileInfoErr := os.Stat(toFile)
		require.False(t, os.IsNotExist(fileInfoErr))
		require.Equal(t, fromFileInfo.Size(), toFileInfo.Size())
		require.True(t, compareFiles(fromFile, toFile))
	})

	t.Run("limit 10, no offset", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 0, 10)

		require.Nil(t, err)
		toFileInfo, fileInfoErr := os.Stat(toFile)
		require.False(t, os.IsNotExist(fileInfoErr))
		require.Equal(t, int64(10), toFileInfo.Size())
		require.True(t, compareFiles("./testdata/out_offset0_limit10.txt", toFile))
	})

	t.Run("limit 1000, no offset", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 0, 1000)

		require.Nil(t, err)
		toFileInfo, fileInfoErr := os.Stat(toFile)
		require.False(t, os.IsNotExist(fileInfoErr))
		require.Equal(t, int64(1000), toFileInfo.Size())
		require.True(t, compareFiles("./testdata/out_offset0_limit1000.txt", toFile))
	})

	t.Run("limit 10000, no offset", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 0, 10000)

		require.Nil(t, err)
		toFileInfo, fileInfoErr := os.Stat(toFile)
		require.False(t, os.IsNotExist(fileInfoErr))
		require.Equal(t, fromFileInfo.Size(), toFileInfo.Size())
		require.True(t, compareFiles("./testdata/out_offset0_limit10000.txt", toFile))
	})

	t.Run("limit 1000, offset 100", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 100, 1000)

		require.Nil(t, err)
		toFileInfo, fileInfoErr := os.Stat(toFile)
		require.False(t, os.IsNotExist(fileInfoErr))
		require.Equal(t, int64(1000), toFileInfo.Size())
		require.True(t, compareFiles("./testdata/out_offset100_limit1000.txt", toFile))
	})

	t.Run("limit 1000, offset 6000", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 6000, 1000)

		require.Nil(t, err)
		require.True(t, compareFiles("./testdata/out_offset6000_limit1000.txt", toFile))
	})

	t.Run("invalid limit", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, 0, -1)

		require.NotNil(t, err)
		require.Truef(t, errors.Is(err, ErrInvalidParameter), "some parameters is invalid", err)
		_, fileInfoErr := os.Stat(toFile)
		require.True(t, os.IsNotExist(fileInfoErr))
	})

	t.Run("invalid offset", func(t *testing.T) {
		os.Remove(toFile)

		err := Copy(fromFile, toFile, -1, 0)

		require.NotNil(t, err)
		require.Truef(t, errors.Is(err, ErrInvalidParameter), "some parameters is invalid", err)
		_, fileInfoErr := os.Stat(toFile)
		require.True(t, os.IsNotExist(fileInfoErr))
	})
}
