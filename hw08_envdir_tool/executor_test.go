package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run(
		"run SIMPLE_BACKUP_SUFFIX=@ VERSION_CONTROL=simple cp ./testdata/env/ADDED ./testdata/COPIED --backup",
		func(t *testing.T) {
			copyFrom := "./testdata/env/BAR"
			copyTo := "./testdata/COPIED"
			backupPrefix := "@"

			env := make(Environment)
			env["VERSION_CONTROL"] = EnvValue{Value: "simple", NeedRemove: false}
			env["SIMPLE_BACKUP_SUFFIX"] = EnvValue{Value: backupPrefix, NeedRemove: false}

			// Копирование файла
			exitCode := RunCmd([]string{"cp", copyFrom, copyTo, "--backup"}, env)
			require.Equal(t, 0, exitCode)

			// Создание бекапа и копирование файла
			exitCode = RunCmd([]string{"cp", copyFrom, copyTo, "--backup"}, env)
			require.Equal(t, 0, exitCode)

			_, err := os.Stat(copyTo)
			require.Nil(t, err, "copied file exists")

			_, err = os.Stat(copyTo + backupPrefix)
			require.Nil(t, err, "copied backup file exists")

			os.Remove(copyTo)
			os.Remove(copyTo + backupPrefix)
		},
	)

	t.Run(
		"run SIMPLE_BACKUP_SUFFIX=@ VERSION_CONTROL=simple cp ./testdata/env/NOT_EXISTS ./testdata/COPIED --backup",
		func(t *testing.T) {
			copyFrom := "./testdata/env/NOT_EXISTS"
			copyTo := "./testdata/COPIED"
			backupPrefix := "!"

			env := make(Environment)
			env["VERSION_CONTROL"] = EnvValue{Value: "simple", NeedRemove: false}
			env["SIMPLE_BACKUP_SUFFIX"] = EnvValue{Value: backupPrefix, NeedRemove: false}

			// Копирование файла, которого нет
			exitCode := RunCmd([]string{"cp", copyFrom, copyTo, "--backup"}, env)
			require.Equal(t, 1, exitCode)
		},
	)
}
