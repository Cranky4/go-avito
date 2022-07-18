package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	environment := make(Environment)

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		if fileInfo.Size() == 0 {
			environment[fileInfo.Name()] = EnvValue{Value: "", NeedRemove: true}
			continue
		}

		file, err := os.Open(dir + "/" + fileInfo.Name())
		if err != nil {
			return nil, err
		}

		reader := bufio.NewReader(file)
		line, err := reader.ReadBytes('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		value := string(line)
		value = strings.Trim(value, "\t\n ")
		value = strings.ReplaceAll(value, "\x00", "\n")

		if len(value) == 0 {
			environment[fileInfo.Name()] = EnvValue{Value: value, NeedRemove: true}
			continue
		}

		environment[fileInfo.Name()] = EnvValue{Value: value, NeedRemove: false}

		file.Close()
	}

	return environment, nil
}
